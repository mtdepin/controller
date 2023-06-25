package processor

import (
	"controller/pkg/logger"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
)

type CallbackChargeProcessor struct {
	stateMachine *statemachine.StateMachine
}

func (p *CallbackChargeProcessor) Init(machine *statemachine.StateMachine) {
	p.stateMachine = machine
}

func (p *CallbackChargeProcessor) Process(request *param.CallbackChargeRequest) (interface{}, error) {
	event, err := p.generateCallbackChargeEvent(request)
	if err != nil {
		return nil, err
	}

	if err := p.stateMachine.Send(event.OrderId, event); err != nil {
		return nil, err
	}
	status := <-event.Ret

	//计费成功,删除此订单状态机
	if request.Status == param.SUCCESS && status == param.SUCCESS {
		if err := p.stateMachine.Delete(request.OrderId); err != nil {
			logger.Warnf("CallbackChargeProcessor  charge finish, delete order: %v statemachine fail: %v ", request.OrderId, err.Error())
		}
	}

	return param.CallbackChargeResponse{
		Status: status,
	}, nil
}

func (p *CallbackChargeProcessor) generateCallbackChargeEvent(request *param.CallbackChargeRequest) (*e.Event, error) {
	callbackDeleteEvent := &e.CallbackChargeEvent{
		OrderId:   request.OrderId,
		OrderType: request.OrderType,
		Status:    request.Status,
	}

	return &e.Event{
		Type:    e.CALLBACK_CHARGE,
		OrderId: request.OrderId,
		Ret:     make(chan int),
		Data:    callbackDeleteEvent,
	}, nil
}
