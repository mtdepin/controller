package statemachine

import (
	"context"
	"controller/pkg/logger"
	sm "controller/pkg/statemachine"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/watcher"
	"encoding/json"
	dss "github.com/ipfs/go-datastore/sync"

	//sm "github.com/filecoin-project/go-statemachine"
	"github.com/ipfs/go-datastore"
)

type StateMachine struct {
	stateGroup *sm.StateGroup
	stateCtl   *StateCtl
}

func (p *StateMachine) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	ds := dss.MutexWrap(datastore.NewMapDatastore())
	p.stateCtl = &StateCtl{}
	p.stateCtl.Init(orderIndex, orderStateIndex)

	p.stateGroup = sm.New(ds, p.stateCtl, OrderState{})
}

func (p *StateMachine) Send(orderId string, event *e.Event) error {
	key := &Key{
		orderId: orderId,
	}

	if err := p.stateGroup.Send(key, event); err != nil {
		bt, _ := json.Marshal(event.Data)
		logger.Warnf("stateGroup send event to statemachine fail, orderId: %v, event: %v, err: %v", orderId, string(bt), err.Error())
		return err
	}
	return nil
}

func (p *StateMachine) Check(orderId string, status uint64) (bool, error) {
	var OrderState OrderState
	stored := p.stateGroup.Get(orderId)
	if err := stored.Get(&OrderState); err != nil {
		return false, err
	}
	return OrderState.Status == status, nil
}

func (p *StateMachine) GetOrderStateInfo(orderId string) (*dict.OrderStateInfo, error) {
	return p.stateCtl.GetOrderStateInfo(orderId)
}

func (p *StateMachine) GetOrderByRequestId(requestId string) (*dict.OrderInfo, error) {
	return p.stateCtl.GetOrderByRequestId(requestId)
}

func (p *StateMachine) GetAllUploadFinishOrderInfo() []*dict.UploadFinishOrder {
	return p.stateCtl.GetAllUploadFinishOrderInfo()
}

func (p *StateMachine) GetBeginReplicateOrderInfo(orderId string) (*dict.UploadFinishOrder, error) {
	return p.stateCtl.GetBeginReplicateOrderInfo(orderId)
}

func (p *StateMachine) UpdateOrderRepInfo(orderId string, tasks map[string]*dict.TaskRepInfo) error {
	return p.stateCtl.UpdateOrderRepInfo(orderId, tasks)
}

func (p *StateMachine) GetOrderStatus(orderId string) (int, error) {
	return p.stateCtl.GetOrderStatus(orderId)
}

func (p *StateMachine) Delete(orderId string) error {
	if err := p.stateCtl.DeleleOrder(orderId); err != nil {
		return err
	}

	//删除前,设置订单状态到prometheus
	watcher.GlobalTraceWatcher.Delete(orderId)
	return p.stateGroup.StopStateMachine(context.Background(), datastore.NewKey(orderId))
}

func (p *StateMachine) UpdateOrder(orderId string, status int) error {
	return p.stateCtl.UpdateOrder(orderId, status)
}

func (p *StateMachine) GetAllOrders(status int) []*dict.OrderStateInfo {
	return p.stateCtl.GetAllOrders(status)
}

func (p *StateMachine) UpdateOrderStateInfo(orderId string, state *dict.OrderStateInfo) error {
	return p.stateCtl.UpdateOrderStateInfo(orderId, state)
}

func (p *StateMachine) GetOrder(orderId string, status int) (*dict.OrderStateInfo, error) {
	return p.stateCtl.GetOrder(orderId, status)
}

func (p *StateMachine) GetFidStatus(orderId, fid string) (int, error) {
	return p.stateCtl.GetFidStatus(orderId, fid)
}

func (p *StateMachine) GetAllOrderIds(status int) []string {
	return p.stateCtl.GetAllOrderIds(status)
}

func (p *StateMachine) GetAllBeginRepOrderInfo() []*dict.UploadFinishOrder {
	return p.stateCtl.GetAllBeginRepOrderInfo()
}

func (p *StateMachine) GetStateFromDB(orderId string) (*dict.OrderStateInfo, error) {
	return p.stateCtl.GetStateFromDB(orderId)
}

func (p *StateMachine) AddPieceFid(orderId string, tasks []*dict.Task) error {
	return p.stateCtl.AddPieceFid(orderId, tasks)
}

func (p *StateMachine) TaskInitFinish(orderId string) (bool, error) {
	return p.stateCtl.TaskInitFinish(orderId)
}
