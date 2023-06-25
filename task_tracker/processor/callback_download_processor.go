package processor

import (
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
	"errors"
	"strings"
)

type CallbackDownloadProcessor struct {
	stateMachine *statemachine.StateMachine
	fidReplicate *database.FidReplication
}

func (p *CallbackDownloadProcessor) Init(machine *statemachine.StateMachine, fidReplicate *database.FidReplication) {
	p.stateMachine = machine
	p.fidReplicate = fidReplicate
}

func (p *CallbackDownloadProcessor) Process(request *param.CallbackDownloadRequest) (interface{}, error) {
	sz := strings.Split(request.Ext, ",")
	if len(sz) < 1 {
		return nil, errors.New("request param err")
	}

	orderId := sz[0]
	fidStatus, er := p.stateMachine.GetFidStatus(orderId, request.Cid)
	if er != nil {
		return param.CallbackDownloadResponse{
			Status: param.SUCCESS,
		}, nil
	}

	if fidStatus >= dict.TASK_DOWNLOAD_SUC { //下载成功或失败， 都不用在处理.
		return param.CallbackDownloadResponse{
			Status: param.SUCCESS,
		}, nil
	}

	event, err := p.generateCallbackDownloadEvent(orderId, request)
	if err != nil {
		return nil, err
	}

	if err := p.stateMachine.Send(event.OrderId, event); err != nil {
		return nil, err
	}
	status := <-event.Ret

	return param.CallbackDownloadResponse{
		Status: status,
	}, nil
}

func (p *CallbackDownloadProcessor) generateCallbackDownloadEvent(orderId string, request *param.CallbackDownloadRequest) (*e.Event, error) {
	callbackDownloadEvent := &e.CallbackDownloadEvent{
		OrderId: orderId,
		Cid:     request.Cid,
		Region:  request.Region,
		Origins: request.Origins,
		Status:  request.Status,
	}

	return &e.Event{
		Type:    e.CALLBACK_DOWNLOAD,
		OrderId: orderId,
		Ret:     make(chan int),
		Data:    callbackDownloadEvent,
	}, nil
}
