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

type DeleteProcessor struct {
	stateMachine   *statemachine.StateMachine
	orderEventChan chan *e.OrderDeleteEvent
	chanSize       int
}

func (p *DeleteProcessor) Init(machine *statemachine.StateMachine) {
	p.stateMachine = machine

	orders := p.stateMachine.GetAllOrders(dict.TASK_DEL_FAIL) //对于备份失败的订单重新备份

	p.chanSize = len(orders) + EXTEND_SIZE
	p.orderEventChan = make(chan *e.OrderDeleteEvent, p.chanSize)

	p.initDeleteEvent(orders)

	logger.Info("----init delete order num: ", len(orders), " deleteChanSize:", p.chanSize)

	go p.Handle()
	//go p.addDeleteEvent()
}

func (p *DeleteProcessor) Add(orderId string) {
	order, err := p.stateMachine.GetOrder(orderId, dict.TASK_DEL_FAIL)
	if err != nil {
		utils.Log(utils.WARN, "DeleteProcessor stateMachine.GetOrder ", err.Error(), nil)
		return
	}

	p.generateDeleteEvent(order)
}

func (p *DeleteProcessor) Handle() {
	for true {
		time.Sleep(TIME_INTERAL * time.Second)

		//一次取n个,防止消息队列中事件太多，导致等待时间过长。
		count := FACTOR * TIME_INTERAL
		nLen := len(p.orderEventChan)
		if count > nLen {
			count = nLen
		}

		for i := 0; i < count; i++ {
			p.delete(<-p.orderEventChan)
			time.Sleep(INTERNAL * time.Millisecond)
		}
	}
}

func (p *DeleteProcessor) delete(event *e.OrderDeleteEvent) {
	rsp, err := utils.Delete(event.Request)
	if err != nil {
		p.addDeleteEventToCache(event)
		return
	}

	state, er := p.stateMachine.GetOrder(event.Request.OrderId, dict.TASK_DEL_FAIL)
	if er != nil {
		utils.Log(utils.WARN, "DeleteProcessor stateMachine.GetRepFailOrder ", er.Error(), event)
		return
	}

	p.update(rsp, state, event.Request)

	if err := p.stateMachine.UpdateOrderStateInfo(event.Request.OrderId, state); err != nil {
		utils.Log(utils.WARN, "DeleteProcessor delete", err.Error(), state)
	}

	if len(event.Request.Tasks) > 0 { //还有区域没备份成功
		p.addDeleteEventToCache(event)
		return
	}

	utils.Log(utils.WARN, "DeleteProcessor delete success", "", rsp)
}

func (p *DeleteProcessor) update(rsp *param.DeleteOrderResponse, state *dict.OrderStateInfo, request *param.DeleteOrderRequest) {
	for _, task := range rsp.Tasks {
		_, ok := state.Tasks[task.Fid]
		if !ok {
			continue
		}

		if task.Status != param.SUCCESS {
			state.Status = dict.TASK_DEL_FAIL

		} else {
			state.Status = dict.TASK_DEL_SUC
		}

		//删除请求中成功的任务
		if _, ok := request.Tasks[task.Fid]; ok {
			if task.Status == dict.TASK_DEL_SUC {
				delete(request.Tasks, task.Fid)
			}
		}
	}

	if len(request.Tasks) == 0 { //全部成功了
		for _, task := range state.Tasks {
			task.Status = dict.TASK_DEL_SUC
		}
		state.Status = dict.TASK_DEL_SUC
	}
}

func (p *DeleteProcessor) addDeleteEventToCache(event *e.OrderDeleteEvent) {
	event.Count++

	if event.Count < dict.DEL_COUNT {
		size := len(p.orderEventChan)
		if size >= p.chanSize-1 {
			utils.Log(utils.WARN, "DeleteProcessor addDeleteEventToCache", fmt.Sprintf(" add  addDeleteEventToCache size = %v  have fill ", size), event.Request)
		}

		p.orderEventChan <- event
	} else {
		utils.Log(utils.WARN, "DeleteProcessor addDeleteEventToCache", fmt.Sprintf("delete count:%d more than %d", event.Count, dict.DEL_COUNT), event.Request)
	}
}

/*func (p *DeleteProcessor) addDeleteEvent() {
	for true {
		time.Sleep(TIME_INTERAL * time.Second)

		//一次取n个,防止消息队列中事件太多，导致等待时间过长。
		count := FACTOR * TIME_INTERAL
		nLen := len(p.orderEventBufChan)
		if count > nLen {
			count = nLen
		}

		for i := 0; i < count; i++ {
			event := <-p.orderEventBufChan
			p.orderEventChan <- event
			time.Sleep(INTERNAL * time.Millisecond)
		}
	}
}*/

func (p *DeleteProcessor) initDeleteEvent(orders []*dict.OrderStateInfo) {
	for _, order := range orders {
		p.generateDeleteEvent(order)
	}
}

func (p *DeleteProcessor) generateDeleteEvent(order *dict.OrderStateInfo) {
	event := &e.OrderDeleteEvent{
		Count: 0,
		Request: &param.DeleteOrderRequest{
			OrderId: order.OrderId,
			Tasks:   make(map[string]*param.UploadTask, 10),
		},
	}

	for _, task := range order.Tasks {
		if task.Status == dict.TASK_DEL_FAIL {
			event.Request.Tasks[task.Fid] = &param.UploadTask{
				Fid:     task.Fid,
				Cid:     task.Cid,
				Origins: task.Origins,
				Region:  task.Region,
			}
		}
	}

	p.orderEventChan <- event
}
