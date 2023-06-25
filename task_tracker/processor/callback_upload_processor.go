package processor

import (
	"controller/pkg/logger"
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
	"gopkg.in/mgo.v2/bson"
)

type CallbackUploadProcessor struct {
	stateMachine *statemachine.StateMachine
	fidReplicate *database.FidReplication
}

func (p *CallbackUploadProcessor) Init(machine *statemachine.StateMachine, fidReplicate *database.FidReplication) {
	p.stateMachine = machine
	p.fidReplicate = fidReplicate
}

func (p *CallbackUploadProcessor) Process(request *param.CallbackUploadRequest) (interface{}, error) {
	fidStatus, er := p.stateMachine.GetFidStatus(request.OrderId, request.Fid) //如果文件已经上传成功了，或者fid不存在，直接返回.
	if er != nil {
		return param.CallbackUploadResponse{
			Status: param.SUCCESS,
		}, nil
	}

	if fidStatus >= dict.TASK_UPLOAD_SUC { //上传成功了， 幂等, 任务备份成功或者失败，都不用在处理.
		logger.Infof("CallbackUploadProcessor Process fid repeat, order_id : %v, fid: %v, request: %v ", request.OrderId, request.Fid, *request)
		return param.CallbackUploadResponse{
			Status: param.SUCCESS,
		}, nil
	}

	if err := p.fidReplicate.Update(request.Fid, bson.M{"$set": bson.M{"cid": request.Cid, "region": request.Region, "origins": request.Origins}}); err != nil {
		logger.Infof("CallbackUploadProcessor ProcessfidReplicate.UpdateCid fail: order_id : %v, fid: %v, err: %v ", request.OrderId, request.Fid, err.Error())
		return nil, err
	}

	logger.Infof("CallbackUploadProcessor ProcessfidReplicate UpdateCid: order_id : %v, fid: %v, region: %v", request.OrderId, request.Fid, request.Region)

	event, err := p.generateCallbackUploadEvent(request)
	if err != nil {
		return nil, err
	}

	if err := p.stateMachine.Send(event.OrderId, event); err != nil {
		return nil, err
	}
	status := <-event.Ret

	return param.CallbackUploadResponse{
		Status: status,
	}, nil
}

func (p *CallbackUploadProcessor) generateCallbackUploadEvent(request *param.CallbackUploadRequest) (*e.Event, error) {
	callbackUploadEvent := &e.CallbackUploadEvent{
		OrderId: request.OrderId,
		Fid:     request.Fid,
		Cid:     request.Cid,
		Region:  request.Region,
		Origins: request.Origins,
		Status:  request.Status,
	}

	return &e.Event{
		Type:    e.CALLBACK_UPLOAD,
		OrderId: request.OrderId,
		Ret:     make(chan int),
		Data:    callbackUploadEvent,
	}, nil
}
