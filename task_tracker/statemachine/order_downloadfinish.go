package statemachine

import (
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
)

type OrderDownloadFinish struct {
	orderIndex      *index.OrderIndex
	orderStateIndex *index.OrderStateIndex
}

func (p *OrderDownloadFinish) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex
}

func (p *OrderDownloadFinish) HandleDownloadFinishEvent(event *e.Event) error {
	orderDownloadFinishEvent := event.Data.(*e.OrderDownloadFinishEvent)
	if orderDownloadFinishEvent.Status == param.SUCCESS {
		return p.handleDownloadFinishSucEvent(orderDownloadFinishEvent)
	}

	return p.handleDownloadFinishFailEvent(orderDownloadFinishEvent)
}

func (p *OrderDownloadFinish) handleDownloadFinishSucEvent(event *e.OrderDownloadFinishEvent) error {
	if err := p.orderStateIndex.SetStatus(event.OrderId, dict.TASK_DOWNLOAD_SUC); err != nil {
		return err
	}

	if err := p.orderIndex.UpdateStatus(event.OrderId, dict.TASK_DOWNLOAD_SUC); err != nil {
		return err
	}

	request, err := p.generateChargeRequest(event.OrderId)
	if err != nil {
		return err
	}

	if _, err := utils.Charge(request); err != nil {
		return p.orderStateIndex.SetStatus(event.OrderId, dict.TASK_CHARGE_FAIL)
	}

	return err
}

func (p *OrderDownloadFinish) handleDownloadFinishFailEvent(event *e.OrderDownloadFinishEvent) error {
	if err := p.orderStateIndex.SetStatus(event.OrderId, dict.TASK_DOWNLOAD_FAIL); err != nil {
		return err
	}

	return p.orderIndex.UpdateStatus(event.OrderId, dict.TASK_DOWNLOAD_FAIL)
}

func (p *OrderDownloadFinish) generateChargeRequest(orderId string) (*param.ChargeRequest, error) {
	state, err := p.orderStateIndex.GetState(orderId)
	if err != nil {
		return nil, err
	}

	return generateChargeRequest(orderId, state)
}
