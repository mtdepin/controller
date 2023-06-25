package index

import (
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	"errors"
	"fmt"
	"sync"
	"time"
)

type OrderIndex struct {
	orderMap   map[string]*dict.OrderInfo //key 订单id, v
	requestMap map[string]string          //key 请求id, value: orderId
	store      *database.Order
	rwLock     *sync.RWMutex
}

func (p *OrderIndex) Init(db *database.DataBase) {
	p.rwLock = new(sync.RWMutex)
	p.store = new(database.Order)
	p.store.Init(db)

	states, err := p.store.Load()
	if err != nil {
		panic(fmt.Sprintf("load orderStateInfo fail, %v", err.Error()))
	}

	p.createIndex(*states)
}

func (p *OrderIndex) createIndex(states []dict.OrderInfo) {
	p.orderMap = make(map[string]*dict.OrderInfo, len(states))
	p.requestMap = make(map[string]string, len(states))
	for i, _ := range states {
		p.orderMap[states[i].OrderId] = &states[i]
		p.requestMap[states[i].RequestId] = states[i].OrderId
	}
}

func (p *OrderIndex) Update(requestId, orderId string, order *dict.OrderInfo) error {
	p.rwLock.Lock()
	p.orderMap[orderId] = order
	p.requestMap[requestId] = order.OrderId
	p.rwLock.Unlock()

	return p.store.Update(order)
}

func (p *OrderIndex) GetOrderStatusByRequestId(requestId string) (int, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	if orderId, ok := p.requestMap[requestId]; ok {
		if order, _ := p.orderMap[orderId]; ok {
			return order.Status, nil
		} else {
			return 0, errors.New(fmt.Sprintf("not find the orderInfo  of requestId: %s, orderId: %s", requestId, orderId))
		}
	} else {
		return 0, errors.New(fmt.Sprintf("not find the order of requestId: %s", requestId))
	}
}

func (p *OrderIndex) GetOrderByRequestId(requestId string) (*dict.OrderInfo, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	if orderId, ok := p.requestMap[requestId]; ok {
		if order, _ := p.orderMap[orderId]; ok {
			return p.copy(order), nil
		} else {
			return nil, errors.New(fmt.Sprintf("not find the orderInfo  of requestId: %s, orderId: %s", requestId, orderId))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("not find the order of requestId: %s", requestId))
	}
}

func (p *OrderIndex) UpdateStatus(orderId string, status int) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	order, ok := p.orderMap[orderId]
	if !ok {
		return errors.New(fmt.Sprintf("OrderIndex UpdateStatus fail, orderId： %v not find", orderId))
	}
	order.Status = status
	order.UpdateTime = time.Now().UnixMilli()

	return p.store.Update(order)
}

func (p *OrderIndex) DeleleOrder(orderId string) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()

	if _, ok := p.orderMap[orderId]; !ok {
		return errors.New(fmt.Sprintf("OrderIndex DeleleOrder err, orderId: %v, not find ", orderId))
	}

	delete(p.orderMap, orderId)

	return nil
}

func (p *OrderIndex) copy(order *dict.OrderInfo) *dict.OrderInfo {
	return &dict.OrderInfo{
		OrderId:    order.OrderId,
		RequestId:  order.RequestId,
		Status:     order.Status,
		Desc:       order.Desc,
		CreateTime: order.CreateTime,
		UpdateTime: order.UpdateTime,
	}
}
