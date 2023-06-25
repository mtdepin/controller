package processor

import (
	"controller/pkg/logger"
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
	"controller/task_tracker/utils"
	"fmt"
	"time"
)

type ChargeProcessor struct {
	stateMachine    *statemachine.StateMachine
	orderEventChan  chan *e.OrderChargeEvent
	chargeOrderChan chan string
	chanSize        int
}

func (p *ChargeProcessor) Init(machine *statemachine.StateMachine, chargeOrderChan chan string) {
	p.stateMachine = machine

	p.chargeOrderChan = chargeOrderChan

	//下载，上传订单.
	orders := p.stateMachine.GetAllOrders(dict.TASK_CHARGE_FAIL)   //对于备份失败的订单重新备份
	repSucOrders := p.stateMachine.GetAllOrders(dict.TASK_REP_SUC) //备份成功订单。

	p.chanSize = len(orders) + EXTEND_SIZE
	p.orderEventChan = make(chan *e.OrderChargeEvent, p.chanSize)

	p.initChargeEvent(orders)
	p.initChargeEvent(repSucOrders)

	logger.Info("----init charge order num: ", len(orders), " chanSize:", p.chanSize)

	go p.Handle()
	go p.addOrder()
}

func (p *ChargeProcessor) Add(orderId string) {
	order, err := p.stateMachine.GetOrder(orderId, dict.TASK_CHARGE_FAIL)
	if err != nil {
		utils.Log(utils.WARN, "ChargeProcessor stateMachine.GetOrder ", err.Error(), nil)
		return
	}

	p.generateChargeEvent(order)
}

func (p *ChargeProcessor) Handle() {
	for true {
		time.Sleep(TIME_INTERAL * time.Second)

		//一次取n个,防止消息队列中事件太多，导致等待时间过长。
		count := FACTOR * TIME_INTERAL
		nLen := len(p.orderEventChan)
		if count > nLen {
			count = nLen
		}

		for i := 0; i < count; i++ {
			p.charge(<-p.orderEventChan)
			time.Sleep(INTERNAL * time.Millisecond)
		}
	}
}

func (p *ChargeProcessor) charge(event *e.OrderChargeEvent) {
	rsp, err := utils.Charge(event.Request)
	if err == nil && rsp.Status == param.SUCCESS {
		if err := p.stateMachine.UpdateOrder(event.Request.OrderId, dict.TASK_CHARGE_SUC); err != nil {
			utils.Log(utils.ERROR, "ChargeProcessor charge stateMachine.UpdateOrder", err.Error(), event)
		}

		return
	}

	p.addEventToCache(event)
}

func (p *ChargeProcessor) addEventToCache(event *e.OrderChargeEvent) {
	event.Count++
	if event.Count < dict.CHARGE_COUNT {
		size := len(p.orderEventChan)
		if size >= p.chanSize-1 {
			utils.Log(utils.WARN, "ChargeProcessor addEventToCache", fmt.Sprintf(" add  addEventToCache size = %v  have fill ", size), event.Request)
		}

		p.orderEventChan <- event
	} else {
		utils.Log(utils.ERROR, "ChargeProcessor addEventToCache", fmt.Sprintf("charge count:%d more than %d", event.Count, dict.CHARGE_COUNT), event.Request)
	}
}

func (p *ChargeProcessor) addOrder() {
	for true {
		order := <-p.chargeOrderChan
		p.Add(order)
	}
}

func (p *ChargeProcessor) initChargeEvent(orders []*dict.OrderStateInfo) {
	for _, order := range orders {
		p.generateChargeEvent(order)
	}
}

func (p *ChargeProcessor) generateChargeEvent(order *dict.OrderStateInfo) {
	event := &e.OrderChargeEvent{
		Count: 0,
		Request: &param.ChargeRequest{
			OrderId:   order.OrderId,
			OrderType: order.OrderType,
			Tasks:     make([]*dict.Task, 0, len(order.Tasks)),
		},
	}

	for _, task := range order.Tasks {
		event.Request.Tasks = append(event.Request.Tasks, task)
	}

	p.orderEventChan <- event
}
