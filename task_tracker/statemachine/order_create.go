package statemachine

import (
	"controller/pkg/logger"
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
		PieceNum:   createOrderEvent.PieceNum,
		Status:     dict.TASK_INIT,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}

	state := &dict.OrderStateInfo{
		OrderId:    event.OrderId,
		OrderType:  int32(createOrderEvent.OrderType),
		PieceNum:   createOrderEvent.PieceNum,
		Tasks:      make(map[string]*dict.Task, createOrderEvent.PieceNum),
		Status:     dict.TASK_INIT,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}

	//上传请求初始化
	/*for _, fid := range createOrderEvent.Fids {
		state.Tasks[fid.Fid] = &dict.Task{Fid: fid.Fid, Cid: fid.Cid, Status: fid.Status, Repeate: fid.Repeate, Origins: fid.Origins}
	}*/

	//下载请求cids
	for _, cid := range createOrderEvent.Cids {
		state.Tasks[cid] = &dict.Task{Cid: cid, Status: dict.TASK_INIT}
	}

	if err := p.orderIndex.Update(order.RequestId, order.OrderId, order); err != nil {
		return err
	}
	logger.Infof(" create_upload_order, order_id: %v, status: %v, update_time: %v", state.OrderId, state.Status, state.UpdateTime)
	return p.orderStateIndex.Update(event.OrderId, state)
}
