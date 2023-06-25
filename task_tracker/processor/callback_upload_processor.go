package processor

import (
	"controller/pkg/logger"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
)

type CallbackUploadProcessor struct {
	stateMachine *statemachine.StateMachine
}

func (p *CallbackUploadProcessor) Init(machine *statemachine.StateMachine) {
	p.stateMachine = machine
}

func (p *CallbackUploadProcessor) Process(request *param.CallbackUploadRequest) (interface{}, error) {
	/*if _, err := p.stateMachine.GetOrderStatus(request.OrderId); err != nil { //订单不存在直接返回成功.
		return param.CallbackUploadResponse{
			Status: param.SUCCESS,
		}, nil
	}*/

	fidStatus, er := p.stateMachine.GetFidStatus(request.OrderId, request.Fid) //如果文件已经上传成功了，或者fid不存在，直接返回.
	if er != nil {
		return param.CallbackUploadResponse{
			Status: param.SUCCESS,
		}, nil
	}

	if fidStatus >= dict.TASK_UPLOAD_SUC { //幂等， 过滤任务已经上传成功了。
		logger.Info("CallbackUploadProcessor Process fid repeat ", request.OrderId, request.Fid, *request)
		return param.CallbackUploadResponse{
			Status: param.SUCCESS,
		}, nil
	}

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
