package processor

import (
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
)

type CallbackRepProcessor struct {
	stateMachine *statemachine.StateMachine
}

func (p *CallbackRepProcessor) Init(machine *statemachine.StateMachine) {
	p.stateMachine = machine
}

func (p *CallbackRepProcessor) Process(request *param.CallbackRepRequest) (interface{}, error) {
	event, err := p.generateCallbackRepEvent(request)
	if err != nil {
		return nil, err
	}

	if err := p.stateMachine.Send(event.OrderId, event); err != nil {
		return nil, err
	}
	status := <-event.Ret

	return param.CallbackRepResponse{
		Status: status,
	}, nil
}

func (p *CallbackRepProcessor) generateCallbackRepEvent(request *param.CallbackRepRequest) (*e.Event, error) {
	callbackRepEvent := &e.CallbackRepEvent{
		OrderId: request.OrderId,
		Fid:     request.Fid,
		Cid:     request.Cid,
		Region:  request.Region,
		Status:  request.Status,
	}

	return &e.Event{
		Type:    e.CALLBACK_REP,
		OrderId: request.OrderId,
		Ret:     make(chan int),
		Data:    callbackRepEvent,
	}, nil
}
