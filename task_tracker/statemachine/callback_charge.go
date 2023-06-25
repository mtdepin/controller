package statemachine

import (
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
)

type CallbackCharge struct {
	orderIndex      *index.OrderIndex
	orderStateIndex *index.OrderStateIndex
}

func (p *CallbackCharge) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex
}

func (p *CallbackCharge) HandleChargeEvent(event *e.Event) error {
	callbackChargeEvent := event.Data.(*e.CallbackChargeEvent)
	if callbackChargeEvent.Status == param.SUCCESS {
		return p.handleChargeSucEvent(callbackChargeEvent)
	}

	return p.handleChargeFailEvent(callbackChargeEvent)
}

func (p *CallbackCharge) handleChargeSucEvent(event *e.CallbackChargeEvent) error {
	if err := p.orderIndex.UpdateStatus(event.OrderId, dict.TASK_CHARGE_SUC); err != nil {
		return err
	}

	return p.orderStateIndex.SetStatus(event.OrderId, dict.TASK_CHARGE_SUC)
}

func (p *CallbackCharge) handleChargeFailEvent(event *e.CallbackChargeEvent) error {
	request, err := p.generateChargeRequest(event.OrderId)
	if err != nil {
		return err
	}

	if _, err := utils.Charge(request); err != nil {
		return p.orderStateIndex.SetStatus(event.OrderId, dict.TASK_CHARGE_FAIL)
	}

	return err
}

func (p *CallbackCharge) generateChargeRequest(orderId string) (*param.ChargeRequest, error) {
	state, err := p.orderStateIndex.GetState(orderId)
	if err != nil {
		return nil, err
	}

	return generateChargeRequest(orderId, state)
}
