package statemachine

import (
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
	"errors"
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

	if orderUploadFinishEvent.Status == param.SUCCESS {
		return p.handleUploadFinishSucEvent(orderUploadFinishEvent)
	}

	return p.handleUploadFinishFailEvent(orderUploadFinishEvent)
}

//备份文件
func (p *OrderUploadFinish) handleUploadFinishSucEvent(event *e.OrderUploadFinishEvent) error {
	strategy, err := utils.GetReplicationStrategy(event.OrderId)
	if err != nil {
		return err
	}

	state, err1 := p.orderStateIndex.GetState(event.OrderId)
	if err1 != nil {
		return err1
	}

	if err := p.createOrderUploadFinishState(strategy, state); err != nil {
		return err
	}

	if err := p.orderStateIndex.Update(event.OrderId, state); err != nil {
		return err
	}

	if err := p.orderIndex.UpdateStatus(event.OrderId, dict.TASK_UPLOAD_SUC); err != nil {
		return err
	}

	if rsp, err := p.replicate(&param.ReplicationRequest{OrderId: event.OrderId, Tasks: state.Tasks}); err != nil {
		SetOrderState(state, dict.TASK_REP_FAIL)
		utils.Log(utils.WARN, "handleUploadFinishSucEvent replicate ", err.Error(), &param.ReplicationRequest{OrderId: event.OrderId, Tasks: state.Tasks})
	} else {
		p.setOrderState(rsp, state)
	}

	return p.orderStateIndex.Update(event.OrderId, state)
}

func (p *OrderUploadFinish) setTaskCid(strategy *param.StrategyInfo, state *dict.OrderStateInfo) {
	for _, val := range strategy.Tasks {
		if task, ok := state.Tasks[val.Fid]; ok {
			val.Cid = task.Cid
		}
	}
}

func (p *OrderUploadFinish) createOrderUploadFinishState(strategy *param.StrategyInfo, state *dict.OrderStateInfo) error {
	if strategy == nil || state == nil {
		return errors.New("param error strategy or state is nil")
	}

	for _, val := range strategy.Tasks {
		if task, ok := state.Tasks[val.Fid]; ok {
			task.Reps = make(map[string]*dict.Rep, len(val.Reps))
			for _, rep := range val.Reps {
				task.Reps[rep.Region] = &dict.Rep{
					Region:     rep.Region,
					VirtualRep: rep.VirtualRep,
					RealRep:    rep.RealRep,
					MinRep:     rep.MinRep,
					MaxRep:     rep.MaxRep,
					Expire:     rep.Expire,
					Encryption: rep.Encryption,
					Status:     dict.TASK_UPLOAD_SUC,
				}
			}
			task.Status = dict.TASK_UPLOAD_SUC
		}
	}

	state.Status = dict.TASK_UPLOAD_SUC //订单上传完成
	state.UpdateTime = time.Now().UnixMilli()

	return nil
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

	for _, task := range state.Tasks { //order 上传失败， 删除已经上传成功的任务。
		if task.Status != dict.TASK_UPLOAD_SUC {
			continue
		}

		request.Tasks[task.Fid] = &param.UploadTask{
			Fid:     task.Fid,
			Cid:     task.Cid,
			Region:  task.Region,
			Origins: task.Origins,
		}
	}

	if len(request.Tasks) > 0 { //有删除任务则删除
		rsp, err := utils.Delete(request)
		if err == nil {
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
		} else { //订单更新成失败
			p.setState(state, dict.TASK_DEL_FAIL)
		}

		if state.Status == dict.TASK_DEL_FAIL {
			utils.Log(utils.WARN, "statemachine handleUploadFinishFailEvent", "delete task fail", rsp)
		}
	} else { //更新订单状为删除成功
		p.setState(state, dict.TASK_DEL_SUC)
	}

	//更新订单表
	if err := p.orderIndex.UpdateStatus(request.OrderId, dict.TASK_DEL_SUC); err != nil {
		return err
	}

	return p.orderStateIndex.Update(request.OrderId, state)
}

//执行批量备份
func (p *OrderUploadFinish) replicate(request *param.ReplicationRequest) (rsp *param.ReplicationResponse, err error) {
	return utils.Replicate(request)
}

func (p *OrderUploadFinish) setOrderState(rsp *param.ReplicationResponse, state *dict.OrderStateInfo) {
	for _, task := range rsp.Tasks {
		stateTask, ok := state.Tasks[task.Fid]
		if !ok {
			continue
		}

		for region, status := range task.RegionStatus {
			if status != param.SUCCESS {
				if rep, ok := stateTask.Reps[region]; ok {
					rep.Status = dict.TASK_REP_FAIL
					state.Status = dict.TASK_REP_FAIL //订单中有一个区域没备份成功，则订单备份失败.
				}
			}
		}
	}
}

func (p *OrderUploadFinish) setState(state *dict.OrderStateInfo, status int) {
	for _, task := range state.Tasks {
		task.Status = status
	}
	state.Status = status
}
