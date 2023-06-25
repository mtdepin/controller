package statemachine

import (
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
	"errors"
	"fmt"
)

type CallbackReplication struct {
	orderIndex      *index.OrderIndex
	orderStateIndex *index.OrderStateIndex
}

func (p *CallbackReplication) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex
}

func (p *CallbackReplication) HandleCallbackRepEvent(event *e.Event) error {
	callbackRepEvent := event.Data.(*e.CallbackRepEvent)

	if callbackRepEvent.Status == param.SUCCESS {
		return p.handleCallbackRepSucEvent(callbackRepEvent)
	}

	return p.handleCallbackRepFailEvent(callbackRepEvent)
}

func (p *CallbackReplication) handleCallbackRepSucEvent(event *e.CallbackRepEvent) error {
	if err := p.orderStateIndex.SetTaskStatus(event.OrderId, event.Fid, event.Region, dict.TASK_REP_SUC); err != nil {
		return err
	}

	if orderStatus, err := p.orderStateIndex.GetOrderStatus(event.OrderId); err == nil && orderStatus == dict.TASK_REP_SUC {
		if err := p.orderIndex.UpdateStatus(event.OrderId, dict.TASK_REP_SUC); err != nil {
			return err
		}

		return p.charge(event.OrderId)
	} else {
		return err
	}
}

func (p *CallbackReplication) charge(orderId string) error {
	request, err := p.generateChargeRequest(orderId)
	if err != nil {
		return err
	}

	if _, err := utils.Charge(request); err != nil {
		return p.orderStateIndex.SetStatus(orderId, dict.TASK_CHARGE_FAIL)
	}

	return err
}

//发送备份命令，进行备份
func (p *CallbackReplication) handleCallbackRepFailEvent(event *e.CallbackRepEvent) error {
	req, err := p.generateRepRequest(event)
	if err != nil {
		return err
	}

	_, err = utils.Replicate(req)

	//备份成功，更新订单状态，返回备份成功事件
	return err
}

func (p *CallbackReplication) generateRepRequest(event *e.CallbackRepEvent) (*param.ReplicationRequest, error) {
	request := &param.ReplicationRequest{
		OrderId: event.OrderId,
		Tasks:   make(map[string]*dict.Task, 1),
	}

	state, err := p.orderStateIndex.GetState(event.OrderId)
	if err != nil {
		return nil, err
	}

	if task, ok := state.Tasks[event.Fid]; ok {
		//to do set one region
		request.Tasks[task.Fid] = task
		return request, nil
	}

	return nil, errors.New(fmt.Sprintf("generateRepRequest not find fid : %v in orderStateIndex", event.Fid))
}

func (p *CallbackReplication) generateChargeRequest(orderId string) (*param.ChargeRequest, error) {
	state, err := p.orderStateIndex.GetState(orderId)
	if err != nil {
		return nil, err
	}

	return generateChargeRequest(orderId, state)
}
