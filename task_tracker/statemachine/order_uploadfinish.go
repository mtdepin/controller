package statemachine

import (
	"controller/pkg/logger"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
	"errors"
	"fmt"
	"time"
)

type OrderUploadFinish struct {
	orderIndex      *index.OrderIndex
	orderStateIndex *index.OrderStateIndex
}

func (p *OrderUploadFinish) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex
}

func (p *OrderUploadFinish) HandleUploadFinishEvent(event *e.Event) error {
	orderUploadFinishEvent := event.Data.(*e.OrderUploadFinishEvent)

	if orderUploadFinishEvent.Status == param.SUCCESS { //开始备份文件
		return p.handleUploadFinishSucEvent(orderUploadFinishEvent)
	}

	return p.handleUploadFinishFailEvent(orderUploadFinishEvent)
}

func (p *OrderUploadFinish) handleUploadFinishSucEvent(event *e.OrderUploadFinishEvent) error {
	status, err := p.orderStateIndex.GetOrderStatus(event.OrderId) //1个文件且存在，通过触发通知发送。
	if err != nil {
		return err
	}

	if status != dict.TASK_UPLOAD_SUC {
		return errors.New(fmt.Sprintf("order: %v status:%v  != dict.TASK_UPLOAD_SUC ", event.OrderId, status))
	}

	return p.replicate(event.OrderId)
}

func (p *OrderUploadFinish) replicate(orderId string) error {
	t1 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish : orderId: %v, GetReplicationStrategy begin: %v ", orderId, t1)
	strategy, err := utils.GetReplicationStrategy(orderId)
	if err != nil {
		return err
	}

	t2 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish : orderId: %v, GetReplicationStrategy end: %v ", orderId, t2-t1)

	state, err := p.orderStateIndex.GetState(orderId)
	if err != nil {
		return err
	}

	if err := setOrderReplicateInfo(strategy, state); err != nil {
		return err
	}

	t3 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish : orderId: %v, Replicate begin: %v ", orderId, t3-t2)

	if rsp, err := utils.Replicate(&param.ReplicationRequest{OrderId: orderId, Tasks: state.Tasks}); err != nil {
		utils.Log(utils.WARN, "handleUploadFinishSucEvent replicate ", err.Error(), &param.ReplicationRequest{OrderId: orderId, Tasks: state.Tasks})
	} else { //根据响应，设置哪些任务备份成功，哪些任务备份失败.state.Tasks = {map[string]*dict.Task}
		setOrderRspState(rsp, state) //添加一个开始备份状态
	}
	//再次更新订单是否开始备份.

	t4 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish : orderId: %v, Replicate end: %v ", orderId, t4-t3)

	if err := p.orderStateIndex.Update(orderId, state); err != nil {
		return err
	}

	t5 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish : orderId: %v, Replicate end: %v ", orderId, t5-t1)
	return p.orderIndex.UpdateStatus(orderId, dict.TASK_BEGIN_REP)
}

//删除文件
func (p *OrderUploadFinish) handleUploadFinishFailEvent(event *e.OrderUploadFinishEvent) error {
	if err := p.orderIndex.UpdateStatus(event.OrderId, dict.TASK_UPLOAD_FAIL); err != nil {
		return err
	}

	state, err := p.orderStateIndex.GetState(event.OrderId)
	if err != nil {
		return err
	}

	request := &param.DeleteOrderRequest{
		OrderId: event.OrderId,
		Tasks:   make(map[string]*param.UploadTask, len(state.Tasks)),
	}

	strategy, er := utils.GetOrderDeleteStrategy(request.OrderId)
	if er != nil {
		return errors.New(fmt.Sprintf("GetDeleteStrategy fail: %v", er.Error()))
	}
	//strategy.Tasks

	for _, v := range strategy.Tasks {
		if task, ok := state.Tasks[v.Fid]; ok { //
			if v.Status == dict.TASK_DEL_SUC {
				task.Status = dict.TASK_DEL_SUC
				continue
			}

			if rep, ok := v.Reps[task.Region]; ok { //策略中，此region, 有订单已关联了此fid, 所以此fid 不能删除, fid直接标记删除成功。
				if rep.MaxRep != 0 || rep.Status == dict.TASK_DEL_SUC {
					task.Status = dict.TASK_DEL_SUC
				}
			} else {
				if task.Region == "" {
					task.Status = dict.TASK_DEL_SUC
				} else {
					utils.Log(utils.WARN, "statemachine handleUploadFinishFailEvent", fmt.Sprintf("region : %v not find in strategy, orderId: %v", task.Region, state.OrderId), strategy)
				}
			}

			if task.Status != dict.TASK_DEL_SUC { //单个集群 删除，不影响其它集群
				request.Tasks[task.Fid] = &param.UploadTask{
					Fid:     task.Fid,
					Cid:     task.Cid,
					Region:  task.Region,
					Origins: task.Origins,
					Status:  dict.TASK_DEL_FAIL, //默认删除失败。
				}
			}
		}
	}

	if len(request.Tasks) > 0 { //有删除任务则删除
		if rsp, err := utils.Delete(request); err == nil {
			for _, task := range rsp.Tasks {
				stateTask, ok := state.Tasks[task.Fid]
				if !ok {
					continue
				}
				if task.Status == param.SUCCESS {
					stateTask.Status = dict.TASK_DEL_SUC
				} else {
					state.Status = dict.TASK_DEL_FAIL
					stateTask.Status = dict.TASK_DEL_FAIL
				}
			}

			if state.Status == dict.TASK_DEL_FAIL {
				utils.Log(utils.WARN, "statemachine handleUploadFinishFailEvent", "delete task fail", request.Tasks)
			}
		}
	} else { //更新订单状为删除成功
		p.setState(state, dict.TASK_DEL_SUC)
	}

	//to do delete order, order_id,备份成功了，在去写入fid,cid.

	//更新订单表
	if err := p.orderIndex.UpdateStatus(request.OrderId, dict.TASK_DEL_SUC); err != nil {
		return err
	}

	return p.orderStateIndex.Update(request.OrderId, state)
}

func (p *OrderUploadFinish) setState(state *dict.OrderStateInfo, status int) {
	for _, task := range state.Tasks {
		task.Status = status
	}
	state.Status = status
}
