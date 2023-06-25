package statemachine

import (
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"time"
)

type OrderCreate struct {
	orderIndex      *index.OrderIndex
	orderStateIndex *index.OrderStateIndex
}

func (p *OrderCreate) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex
}

func (p *OrderCreate) HandleCreateOrderEvent(event *e.Event) error {
	createOrderEvent := event.Data.(*e.CreateOrderEvent)
	order := &dict.OrderInfo{
		OrderId:    event.OrderId,
		RequestId:  createOrderEvent.RequestId,
		OrderType:  createOrderEvent.OrderType,
		Status:     dict.TASK_INIT,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}

	state := &dict.OrderStateInfo{
		OrderId:    event.OrderId,
		OrderType:  int32(createOrderEvent.OrderType),
		Status:     dict.TASK_INIT,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}

	state.Tasks = make(map[string]*dict.Task, 10)

	for _, fid := range createOrderEvent.Fids {
		state.Tasks[fid] = &dict.Task{Fid: fid}
	}

	for _, cid := range createOrderEvent.Cids {
		state.Tasks[cid] = &dict.Task{Cid: cid}
	}

	if err := p.orderIndex.Update(order.RequestId, order.OrderId, order); err != nil {
		return err
	}

	return p.orderStateIndex.Update(event.OrderId, state)
}
