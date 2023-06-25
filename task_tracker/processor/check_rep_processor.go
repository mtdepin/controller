package processor

import (
	"controller/pkg/logger"
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
	"controller/task_tracker/utils"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type CheckRepProcessor struct {
	stateMachine           *statemachine.StateMachine
	orderRepCheckEventChan chan *e.OrderRepCheckEvent
	fidReplicate           *database.FidReplication
	repOrderChan           chan string
	chargeOrderChan        chan string
	repChanSize            int //extend
}

func (p *CheckRepProcessor) Init(machine *statemachine.StateMachine, repOrderChan, chargeOrderChan chan string, fidReplicate *database.FidReplication) {
	p.stateMachine = machine
	p.fidReplicate = fidReplicate
	p.repOrderChan = repOrderChan
	p.chargeOrderChan = chargeOrderChan
	orders := p.stateMachine.GetAllBeginRepOrderInfo()

	p.repChanSize = len(orders) + EXTEND_SIZE
	p.orderRepCheckEventChan = make(chan *e.OrderRepCheckEvent, p.repChanSize)

	p.initRepCheckEvent(orders)
	logger.Info("----init search replicate order num: ", len(orders), " repChanSize:", p.repChanSize)

	go p.Handle()
	go p.addOrder()
}

func (p *CheckRepProcessor) Add(orderId string) {
	order, err := p.stateMachine.GetBeginReplicateOrderInfo(orderId)
	if err != nil {
		utils.Log(utils.WARN, "CheckRepProcessor Add ", err.Error(), nil)
		return
	}
	p.generateRepCheckEvent(order)
}

func (p *CheckRepProcessor) Handle() {
	for true {
		time.Sleep(TIME_INTERAL * time.Second)
		//一次取n个,防止消息队列中事件太多，导致等待时间过长。
		count := FACTOR * TIME_INTERAL
		nLen := len(p.orderRepCheckEventChan)
		if count > nLen {
			count = nLen
		}

		for i := 0; i < count; i++ {
			p.SearchReplicate(<-p.orderRepCheckEventChan)
			time.Sleep(INTERNAL * time.Millisecond)
		}
	}
}

func (p *CheckRepProcessor) SearchReplicate(event *e.OrderRepCheckEvent) {
	//check 订单是否备份成功:
	orderId := event.Request.OrderId
	status, err := p.stateMachine.GetOrderStatus(orderId)
	if err == nil && status == dict.TASK_REP_SUC { //如果订单备份成功，直接返回.
		p.stateMachine.UpdateOrder(orderId, status)

		p.UpdateFidRepStatus(orderId)

		event := p.generateCallbackChargeEvent(orderId, param.FAIL)
		p.stateMachine.Send(orderId, event)
		ret := <-event.Ret

		if ret != param.SUCCESS {
			p.chargeOrderChan <- orderId //重新计费
		}
		return
	}

	rsp, err := utils.SearchRep(event.Request)
	if err != nil { //查询失败，直接重新查询
		p.addSearchEvent(event) //可能block, to do proc
		return
	}

	//查询成功， 更新订单状态， 如果有任务备份失败，投放备份事件；  如果全部备份成功，投放计费事件。对于还没备份完成的任务，重新查询。
	err = p.stateMachine.UpdateOrderRepInfo(rsp.OrderId, rsp.Tasks)
	if err != nil {
		utils.Log(utils.WARN, "CheckRepProcessor.SearchReplicate", err.Error(), rsp)
		if _, err := p.stateMachine.GetOrderStatus(rsp.OrderId); err != nil {
			return //订单不存在，则返回.
		}
		p.addSearchEvent(event) //重新查询
		return
	}

	//生成重新查询请求 及 重新备份事件。
	request, events := p.generate(rsp.OrderId, rsp.Tasks)
	for _, event := range events {
		p.stateMachine.Send(rsp.OrderId, event)
		ret := <-event.Ret

		if ret != param.SUCCESS {
			utils.Log(utils.WARN, "CheckRepProcessor.SearchReplicate, stateMachine proc replicate", "fail", rsp)
		}
	}

	if request != nil { //重新查询未备份完成的请求.
		p.addSearchEvent(&e.OrderRepCheckEvent{Count: event.Count, BeginTime: event.BeginTime, Request: request}) //count 累加.
	} else { //全部备份成功，生成计费事件.
		status, err := p.stateMachine.GetOrderStatus(rsp.OrderId)
		if err == nil && status == dict.TASK_REP_SUC {
			p.stateMachine.UpdateOrder(rsp.OrderId, status)
			//to do set fid replicate success status.

			p.UpdateFidRepStatus(rsp.OrderId)

			event := p.generateCallbackChargeEvent(rsp.OrderId, param.FAIL)
			p.stateMachine.Send(rsp.OrderId, event)
			ret := <-event.Ret

			if ret != param.SUCCESS {
				p.chargeOrderChan <- rsp.OrderId //重新计费
			}
		} else { //重新查询.
			if err == nil {
				p.addSearchEvent(event) //获取订单状态失败，或者等你的没有备份成功，重新查询。
			} else {
				utils.Log(utils.ERROR, "CheckRepProcessor.SearchReplicate, stateMachine.GetOrderStatus fail", err.Error(), event)
			}
		}
	}
}

func (p *CheckRepProcessor) generate(orderId string, tasks map[string]*dict.TaskRepInfo) (*dict.UploadFinishOrder, []*e.Event) {
	request := &dict.UploadFinishOrder{
		OrderId: orderId,
		Tasks:   make([]*dict.RepTask, 0, len(tasks)),
	}

	events := make([]*e.Event, 0, 10)

	for _, task := range tasks {
		taskRequest := &dict.RepTask{
			Fid:     task.Fid,
			Cid:     task.Cid,
			Regions: make([]string, 0, len(task.Regions)),
		}

		for _, region := range task.Regions {
			if region.Status == dict.FAIL { //失败，生成重新备份事件。
				event := p.generateCallbackRepEvent(orderId, task.Fid, task.Cid, region.Region, region.Status)
				events = append(events, event)
			}

			if region.CheckStatus != dict.SUCESS {
				taskRequest.Regions = append(taskRequest.Regions, region.Region)
			}
		}

		if len(taskRequest.Regions) > 0 {
			request.Tasks = append(request.Tasks, taskRequest)
		}
	}

	if len(request.Tasks) > 0 {
		return request, events
	}

	return nil, events
}

func (p *CheckRepProcessor) addSearchEvent(event *e.OrderRepCheckEvent) {
	event.Count++
	diffTime := time.Now().UnixMilli() - event.BeginTime
	if event.Count < dict.SEARCH_COUNT && diffTime < dict.Duration {
		size := len(p.orderRepCheckEventChan)
		if size >= p.repChanSize-1 {
			utils.Log(utils.ERROR, "CheckRepProcessor addSearchEvent ", fmt.Sprintf(" add  orderRepCheckEventChan size = %v  have fill ", size), event.Request)
		}

		p.orderRepCheckEventChan <- event
	} else {
		utils.Log(utils.ERROR, "CheckRepProcessor addSearchEvent ", fmt.Sprintf("seach more than count: %d  or more than time: %d ", dict.SEARCH_COUNT, diffTime), event.Request)
	}
}

func (p *CheckRepProcessor) initRepCheckEvent(orders []*dict.UploadFinishOrder) {
	for _, order := range orders {
		p.generateRepCheckEvent(order)
	}
}

func (p *CheckRepProcessor) generateRepCheckEvent(order *dict.UploadFinishOrder) {
	size := len(p.orderRepCheckEventChan)

	if size >= p.repChanSize-1 {
		utils.Log(utils.ERROR, "CheckRepProcessor, generateRepCheckEvent", fmt.Sprintf("add event to buffer, orderRepCheckEventChan buffer have fill, size= %v", size), order)
	}

	p.orderRepCheckEventChan <- &e.OrderRepCheckEvent{
		Count:     0,
		BeginTime: time.Now().UnixMilli(), //ms
		Request:   order,
	}
}

func (p *CheckRepProcessor) generateCallbackRepEvent(orderId, fid, cid, region string, status int) *e.Event {
	callbackRepEvent := &e.CallbackRepEvent{
		OrderId: orderId,
		Fid:     fid,
		Cid:     cid,
		Region:  region,
		Status:  status,
	}

	return &e.Event{
		Type:    e.CALLBACK_REP,
		OrderId: orderId,
		Ret:     make(chan int),
		Data:    callbackRepEvent,
	}
}

func (p *CheckRepProcessor) generateCallbackChargeEvent(orderId string, status int) *e.Event {
	callbackChargeEvent := &e.CallbackChargeEvent{
		OrderId:   orderId,
		OrderType: param.UPLOAD,
		Status:    status,
	}

	return &e.Event{
		Type:    e.CALLBACK_CHARGE,
		OrderId: orderId,
		Ret:     make(chan int),
		Data:    callbackChargeEvent,
	}
}

func (p *CheckRepProcessor) addOrder() {
	for true {
		order := <-p.repOrderChan
		p.Add(order)
	}
}

func (p *CheckRepProcessor) UpdateFidRepStatus(orderId string) error {
	state, err := p.stateMachine.GetOrderStateInfo(orderId)
	if err != nil {
		return err
	}

	for _, task := range state.Tasks {
		if task.Status == dict.TASK_REP_SUC {
			if err := p.fidReplicate.Update(task.Fid, bson.M{"$set": bson.M{"status": task.Status}}); err != nil {
				logger.Warnf("CheckRepProcessor UUpdateFidRepStatus fail order_id : %v, fid: %v, err: %v ", orderId, task.Fid, err.Error())
			}
		}
	}

	return nil
}
