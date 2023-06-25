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

	orders := p.stateMachine.GetAllOrderIds(dict.TASK_DEL_FAIL)

	p.chanSize = len(orders) + EXTEND_SIZE
	p.orderEventChan = make(chan *e.OrderDeleteEvent, p.chanSize)

	p.initDeleteEvent(orders)

	logger.Info("----init delete order num: ", len(orders), " deleteChanSize:", p.chanSize)

	go p.Handle()
	//go p.addDeleteEvent()
}

func (p *DeleteProcessor) Add(orderId string) {
	p.addDeleteEventToCache(&e.OrderDeleteEvent{Count: 0, OrderId: orderId})
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
	//获取删除策略，生成删除请求。
	state, err := p.stateMachine.GetOrder(event.OrderId, dict.TASK_DEL_FAIL)
	if err != nil {
		utils.Log(utils.WARN, "DeleteProcessor delete stateMachine.GetOrder  fail", err.Error(), event)
		return
	}

	strategy, err := utils.GetOrderDeleteStrategy(event.OrderId)
	if err != nil {
		utils.Log(utils.WARN, "DeleteProcessor delete  GetReplicationStrategy ", err.Error(), event)
		return
	}

	deleteRequest := p.generateDeleteRequest(state, strategy)
	if state.Status == dict.TASK_DEL_SUC {
		if err := p.stateMachine.UpdateOrderStateInfo(event.OrderId, state); err != nil {
			utils.Log(utils.WARN, "DeleteProcessor delete stateMachine.UpdateOrderStateInfo fail ", err.Error(), state)
		}
		return
	}

	rsp, err := utils.Delete(deleteRequest)
	if err != nil {
		p.addDeleteEventToCache(event) //重新获取策略。
		return
	}

	p.updateOrderStatus(rsp, state)

	if err := p.stateMachine.UpdateOrderStateInfo(event.OrderId, state); err != nil {
		utils.Log(utils.WARN, "DeleteProcessor delete", err.Error(), state)
	}

	if state.Status != dict.TASK_DEL_SUC { //还有区域没删除成功
		p.addDeleteEventToCache(event)
		return
	}

	utils.Log(utils.WARN, "DeleteProcessor delete success", "", rsp)
}

func (p *DeleteProcessor) updateOrderStatus(rsp *param.DeleteOrderResponse, state *dict.OrderStateInfo) {
	for _, task := range rsp.Tasks {
		stTask, ok := state.Tasks[task.Fid]
		if !ok {
			continue
		}

		if task.Status != param.SUCCESS {
			stTask.Status = dict.TASK_DEL_FAIL
			state.Status = dict.TASK_DEL_FAIL
		} else {
			stTask.Status = dict.TASK_DEL_SUC
		}
	}
}

func (p *DeleteProcessor) addDeleteEventToCache(event *e.OrderDeleteEvent) {
	event.Count++

	if event.Count < dict.DEL_COUNT {
		size := len(p.orderEventChan)
		if size >= p.chanSize-1 {
			utils.Log(utils.ERROR, "DeleteProcessor addDeleteEventToCache", fmt.Sprintf(" add  addDeleteEventToCache size = %v  have fill ", size), event)
		}

		p.orderEventChan <- event
	} else {
		utils.Log(utils.ERROR, "DeleteProcessor addDeleteEventToCache", fmt.Sprintf("delete count:%d more than %d", event.Count, dict.DEL_COUNT), event)
	}
}

func (p *DeleteProcessor) initDeleteEvent(orderIds []string) {
	for _, orderId := range orderIds {
		p.addDeleteEventToCache(&e.OrderDeleteEvent{Count: 0, OrderId: orderId})
	}
}

func (p *DeleteProcessor) generateDeleteRequest(order *dict.OrderStateInfo, strategy *param.StrategyInfo) *param.DeleteOrderRequest {
	request := &param.DeleteOrderRequest{
		OrderId: order.OrderId,
		Tasks:   make(map[string]*param.UploadTask, 10),
	}

	//遍历策略,看是否删除
	//to do proc
	for _, stTask := range strategy.Tasks {
		if task, ok := order.Tasks[stTask.Fid]; ok {
			if task.Status == dict.TASK_DEL_FAIL { //如果任务删除失败。
				if rep, ok := stTask.Reps[task.Region]; ok {
					if rep.Status == dict.TASK_DEL_SUC || rep.MaxRep > 0 { //有cid 关联，不能删除.
						task.Status = dict.TASK_DEL_SUC
					}
				} else {
					if task.Region == "" { //如果region 为空,则设置成删除成功。
						task.Status = dict.TASK_DEL_SUC
					} else {
						utils.Log(utils.WARN, "DeleteProcessor generateDeleteRequest", fmt.Sprintf("delete task region %v not find in strategy, orderId:%v", task.Region, order.OrderId), strategy)
					}
				}

				if task.Status != dict.TASK_DEL_SUC {
					request.Tasks[task.Fid] = &param.UploadTask{
						Fid:     task.Fid,
						Cid:     task.Cid,
						Origins: task.Origins,
						Region:  task.Region,
					}
				}
			}
		}
	}

	if len(request.Tasks) == 0 {
		order.Status = dict.TASK_DEL_SUC
	}

	return request
}
