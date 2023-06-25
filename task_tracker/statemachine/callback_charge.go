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

	state, err := p.orderStateIndex.GetState(event.OrderId)
	if err != nil {
		return err
	}
	for _, task := range state.Tasks {
		for _, rep := range task.Reps {
			rep.Status = dict.TASK_CHARGE_SUC
		}
		task.Status = dict.TASK_CHARGE_SUC
	}
	state.Status = dict.TASK_CHARGE_SUC

	return p.orderStateIndex.Update(event.OrderId, state)
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
