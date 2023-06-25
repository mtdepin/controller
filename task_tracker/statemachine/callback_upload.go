package statemachine

import (
	"controller/pkg/logger"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
	"time"
)

type CallbackUpload struct {
	orderIndex      *index.OrderIndex
	orderStateIndex *index.OrderStateIndex
}

func (p *CallbackUpload) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex
}

func (p *CallbackUpload) HandleUploadEvent(event *e.Event) error {
	callbackUploadEvent := event.Data.(*e.CallbackUploadEvent)

	if callbackUploadEvent.Status == param.SUCCESS {
		return p.handleUploadSucEvent(callbackUploadEvent)
	}
	return p.handleUploadFailEvent(callbackUploadEvent)
}

func (p *CallbackUpload) handleUploadSucEvent(event *e.CallbackUploadEvent) error {
	if err := p.orderStateIndex.SetTaskUploadStatus(event.OrderId, event.Fid, event.Cid, event.Region, event.Origins, dict.TASK_UPLOAD_SUC); err != nil {
		return err
	}

	//订单不存在，或者订单状态还没上传成功，则返回成功。
	if status, err := p.orderStateIndex.GetOrderStatus(event.OrderId); err != nil || status != dict.TASK_UPLOAD_SUC {
		return nil
	}

	//文件上传成功，开始备份.
	return p.Replicate(event.OrderId)
}

func (p *CallbackUpload) handleUploadFailEvent(event *e.CallbackUploadEvent) error {
	return p.orderStateIndex.SetTaskUploadStatus(event.OrderId, event.Fid, event.Cid, event.Region, event.Origins, dict.TASK_UPLOAD_FAIL)
}

//备份文件
func (p *CallbackUpload) Replicate(orderId string) error {
	t1 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish  CallbackUpload: orderId: %v, GetReplicationStrategy begin: %v ", orderId, t1)

	strategy, err := utils.GetReplicationStrategy(orderId)
	if err != nil {
		return err
	}

	t2 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish  CallbackUpload: orderId: %v, GetReplicationStrategy end: %v ms", orderId, t2-t1)

	state, err := p.orderStateIndex.GetState(orderId)
	if err != nil {
		return err
	}

	if err := setOrderReplicateInfo(strategy, state); err != nil {
		return err
	}

	t3 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish  CallbackUpload: orderId: %v, Replicate begin: %v ms", orderId, t3-t2)
	if rsp, err := utils.Replicate(&param.ReplicationRequest{OrderId: orderId, Tasks: state.Tasks}); err != nil {
		utils.Log(utils.WARN, "handleUploadFinishSucEvent replicate ", err.Error(), &param.ReplicationRequest{OrderId: orderId, Tasks: state.Tasks})
	} else { //根据响应，设置哪些任务备份成功，哪些任务备份失败.state.Tasks = {map[string]*dict.Task}
		setOrderRspState(rsp, state) //添加一个开始备份状态
	}
	//再次更新订单是否开始备份.
	t4 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish  CallbackUpload: orderId: %v, Replicate end: %v ms", orderId, t4-t3)
	if err := p.orderStateIndex.Update(orderId, state); err != nil {
		return err
	}
	t5 := time.Now().UnixMilli()
	logger.Infof(" uploadfinish  CallbackUpload: orderId: %v, total time: %v ms", orderId, t5-t1)
	return p.orderIndex.UpdateStatus(orderId, dict.TASK_BEGIN_REP)
}
