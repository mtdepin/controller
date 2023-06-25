package statemachine

import (
	"controller/pkg/logger"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/utils"
	"encoding/json"
	"fmt"
	//sm "github.com/filecoin-project/go-statemachine"
	sm "controller/pkg/statemachine"
)

type StateCtl struct {
	orderIndex          *index.OrderIndex
	orderStateIndex     *index.OrderStateIndex
	callbackCharge      *CallbackCharge
	callbackUpload      *CallbackUpload
	CallbackReplicate   *CallbackReplication
	orderCreate         *OrderCreate
	orderUploadFinish   *OrderUploadFinish
	orderDownloadFinish *OrderDownloadFinish
}

func (p *StateCtl) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex

	p.callbackCharge = new(CallbackCharge)
	p.callbackCharge.Init(orderIndex, orderStateIndex)

	p.callbackUpload = new(CallbackUpload)
	p.callbackUpload.Init(orderIndex, orderStateIndex)

	p.orderCreate = new(OrderCreate)
	p.orderCreate.Init(orderIndex, orderStateIndex)

	p.orderUploadFinish = new(OrderUploadFinish)
	p.orderUploadFinish.Init(orderIndex, orderStateIndex)

	p.CallbackReplicate = new(CallbackReplication)
	p.CallbackReplicate.Init(orderIndex, orderStateIndex)
	p.orderDownloadFinish = new(OrderDownloadFinish)
	p.orderDownloadFinish.Init(orderIndex, orderStateIndex)
}

func (p *StateCtl) Plan(events []sm.Event, state interface{}) (interface{}, uint64, error) {
	return p.plan(events, state.(*OrderState))
}

func (p *StateCtl) plan(events []sm.Event, orderState *OrderState) (func(sm.Context, OrderState) error, uint64, error) {
	event := events[0]
	orderEvent := event.User.(*e.Event)
	orderState.Event = orderEvent

	switch orderEvent.Type {
	case e.CREATE_ORDER:
		orderState.Status = INIT
	case e.CALLBACK_UPLOAD:
		orderState.Status = UPLOAD_PROCEED
	case e.TASK_UPLOAD_FINISH:
		orderState.Status = UPLOAD_FINISH
	case e.CALLBACK_REP:
		orderState.Status = REP_PROCEED
	case e.CALLBACK_CHARGE:
		orderState.Status = CHARGE_PROCEED
	case e.TASK_DOWNLOAD_FINISH:
		orderState.Status = DOWNLOAD_FINISH
	}

	switch orderState.Status {
	case INIT:
		return p.InitOrder, 1, nil
	case UPLOAD_PROCEED:
		return p.UploadProceed, 1, nil
	case UPLOAD_FINISH:
		return p.UploadFinish, 1, nil
	case REP_PROCEED:
		return p.RepProceed, 1, nil
	case CHARGE_PROCEED:
		return p.ChargeProceed, 1, nil
	case DOWNLOAD_FINISH:
		return p.DownloadFinish, 1, nil
	default:
		bt, _ := json.Marshal(orderState.Event)
		logger.Errorf("statectl plan  error order status: %d, event: %v", orderState.Status, string(bt))
	}

	panic(fmt.Sprintf("statectl plan invalid order status %d", orderState.Status))
}

func (p *StateCtl) InitOrder(ctx sm.Context, st OrderState) error {
	if err := p.orderCreate.HandleCreateOrderEvent(st.Event); err != nil {
		utils.Log(utils.WARN, "StateCtl HandleCreateOrderEvent", err.Error(), st.Event.Data)
		st.Event.Ret <- FAIL
		return nil
	}
	st.Event.Ret <- SUCCESS
	return nil
}

func (p *StateCtl) UploadProceed(ctx sm.Context, st OrderState) error {
	if err := p.callbackUpload.HandleUploadEvent(st.Event); err != nil {
		utils.Log(utils.WARN, "StateCtl HandleUploadEvent", err.Error(), st.Event.Data)
		st.Event.Ret <- FAIL
		return nil
	}
	st.Event.Ret <- SUCCESS
	return nil
}

func (p *StateCtl) UploadFinish(ctx sm.Context, st OrderState) error {
	if err := p.orderUploadFinish.HandleUploadFinishEvent(st.Event); err != nil {
		utils.Log(utils.WARN, "StateCtl HandleUploadFinishEvent", err.Error(), st.Event.Data)
		st.Event.Ret <- FAIL
		return nil
	}
	st.Event.Ret <- SUCCESS
	return nil
}

func (p *StateCtl) RepProceed(ctx sm.Context, st OrderState) error {
	if err := p.CallbackReplicate.HandleCallbackRepEvent(st.Event); err != nil {
		utils.Log(utils.WARN, "StateCtl HandleCallbackRepEvent", err.Error(), st.Event.Data)
		st.Event.Ret <- FAIL
		return nil
	}
	st.Event.Ret <- SUCCESS
	return nil
}

/*func (p *StateCtl) DeleteProceed(ctx sm.Context, st OrderState) error {
	if err := p.callbackDelete.HandleDeleteEvent(st.Event); err != nil {
		utils.Log(utils.WARN, "StateCtl HandleDeleteEvent", err.Error(), st.Event.Data)
		st.Event.Ret <- FAIL
		return nil
	}
	st.Event.Ret <- SUCCESS
	return nil
}*/

func (p *StateCtl) ChargeProceed(ctx sm.Context, st OrderState) error {
	if err := p.callbackCharge.HandleChargeEvent(st.Event); err != nil {
		utils.Log(utils.WARN, "StateCtl HandleChargeEvent", err.Error(), st.Event.Data)
		st.Event.Ret <- FAIL
		return nil
	}
	st.Event.Ret <- SUCCESS
	return nil
}

func (p *StateCtl) DownloadFinish(ctx sm.Context, st OrderState) error {
	if err := p.orderDownloadFinish.HandleDownloadFinishEvent(st.Event); err != nil {
		utils.Log(utils.WARN, "StateCtl HandleDownloadFinishEvent", err.Error(), st.Event.Data)
		st.Event.Ret <- FAIL
		return nil
	}
	st.Event.Ret <- SUCCESS
	return nil
}

func (p *StateCtl) GetOrderStateInfo(orderId string) (*dict.OrderStateInfo, error) {
	return p.orderStateIndex.GetState(orderId)
}

func (p *StateCtl) GetOrderByRequestId(requestId string) (*dict.OrderInfo, error) {
	return p.orderIndex.GetOrderByRequestId(requestId)
}

func (p *StateCtl) GetAllUploadFinishOrderInfo() []*dict.UploadFinishOrder {
	return p.orderStateIndex.GetAllUploadFinishOrderInfo()
}

func (p *StateCtl) GetUploadFinishOrderInfo(orderId string) (*dict.UploadFinishOrder, error) {
	return p.orderStateIndex.GetUploadFinishOrderInfo(orderId)
}

func (p *StateCtl) UpdateOrderRepInfo(orderId string, tasks map[string]*dict.TaskRepInfo) error {
	return p.orderStateIndex.UpdateOrderRepInfo(orderId, tasks)
}

func (p *StateCtl) GetOrderStatus(orderId string) (int, error) {
	return p.orderStateIndex.GetOrderStatus(orderId)
}

func (p *StateCtl) DeleleOrder(orderId string) error {
	if err := p.orderIndex.DeleleOrder(orderId); err != nil {
		return err
	}

	return p.orderStateIndex.DeleleOrder(orderId)
}

func (p *StateCtl) UpdateOrder(orderId string, status int) error {
	return p.orderIndex.UpdateStatus(orderId, status)
}

func (p *StateCtl) GetAllOrders(status int) []*dict.OrderStateInfo {
	return p.orderStateIndex.GetAllOrders(status)
}

func (p *StateCtl) UpdateOrderStateInfo(orderId string, state *dict.OrderStateInfo) error {
	return p.orderStateIndex.Update(orderId, state)
}

func (p *StateCtl) GetOrder(orderId string, status int) (*dict.OrderStateInfo, error) {
	return p.orderStateIndex.GetOrder(orderId, status)
}

func (p *StateCtl) GetFidStatus(orderId, fid string) (int, error) {
	return p.orderStateIndex.GetFidStatus(orderId, fid)
}
