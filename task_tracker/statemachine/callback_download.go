package statemachine

import (
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
)

type CallbackDownload struct {
	orderIndex      *index.OrderIndex
	orderStateIndex *index.OrderStateIndex
}

func (p *CallbackDownload) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex
}

func (p *CallbackDownload) HandleDownloadEvent(event *e.Event) error {
	CallbackDownloadEvent := event.Data.(*e.CallbackDownloadEvent)

	if CallbackDownloadEvent.Status == param.SUCCESS {
		return p.handleDownloadSucEvent(CallbackDownloadEvent)
	}
	return p.handleDownloadFailEvent(CallbackDownloadEvent)
}

func (p *CallbackDownload) handleDownloadSucEvent(event *e.CallbackDownloadEvent) error {
	if err := p.orderStateIndex.SetTaskDownloadStatus(event.OrderId, event.Cid, event.Region, event.Origins, dict.TASK_DOWNLOAD_SUC); err != nil {
		return err
	}

	//订单不存在，或者订单状态还没下载成功，则返回成功。
	if status, err := p.orderStateIndex.GetOrderStatus(event.OrderId); err != nil || status != dict.TASK_DOWNLOAD_SUC {
		return nil
	}

	//订单已经下载完成，开始计费。
	request, err := p.generateChargeRequest(event.OrderId)
	if err != nil {
		return err
	}

	if _, err := utils.Charge(request); err != nil {
		return p.orderStateIndex.SetStatus(event.OrderId, dict.TASK_CHARGE_FAIL)
	}

	return err
}

func (p *CallbackDownload) handleDownloadFailEvent(event *e.CallbackDownloadEvent) error {
	if err := p.orderStateIndex.SetTaskDownloadStatus(event.OrderId, event.Cid, event.Region, event.Origins, dict.TASK_DOWNLOAD_SUC); err != nil {
		return err
	}

	return p.orderIndex.UpdateStatus(event.OrderId, dict.TASK_DOWNLOAD_FAIL)
}

func (p *CallbackDownload) generateChargeRequest(orderId string) (*param.ChargeRequest, error) {
	state, err := p.orderStateIndex.GetState(orderId)
	if err != nil {
		return nil, err
	}

	return generateChargeRequest(orderId, state)
}
