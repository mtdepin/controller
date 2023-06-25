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

type ReplicateProcessor struct {
	stateMachine      *statemachine.StateMachine
	orderRepEventChan chan *e.OrderRepEvent
	orderChan         chan string
	repChanSize       int
}

func (p *ReplicateProcessor) Init(machine *statemachine.StateMachine, orderChan chan string) {
	p.stateMachine = machine
	p.orderChan = orderChan

	orders := p.stateMachine.GetAllOrderIds(dict.TASK_REP_FAIL) //备份失败的订单重新备份

	uploadFinishOrders := p.stateMachine.GetAllOrderIds(dict.TASK_UPLOAD_SUC) //极端情况，先把上传成功文件统一处理成上传失败状态，要从新获取备份策略进行备份。

	p.repChanSize = len(orders) + len(uploadFinishOrders) + EXTEND_SIZE
	p.orderRepEventChan = make(chan *e.OrderRepEvent, p.repChanSize)

	logger.Info("----init replicate order total_num: ", len(orders)+len(uploadFinishOrders), " rep_fail order_num: ", len(orders), " upload finish order_num: ", len(uploadFinishOrders), " repChanSize:", p.repChanSize)
	p.initRepEvent(orders)
	p.processUploadFinishOrder(uploadFinishOrders)

	go p.Handle()
}

func (p *ReplicateProcessor) Add(orderId string) {
	p.addRepEventToCache(&e.OrderRepEvent{Count: 0, OrderId: orderId})
}

func (p *ReplicateProcessor) Handle() {
	for true {
		time.Sleep(TIME_INTERAL * time.Second)
		//一次取n个,防止消息队列中事件太多，导致等待时间过长。
		count := FACTOR * TIME_INTERAL
		nlen := len(p.orderRepEventChan)
		if count > nlen {
			count = nlen
		}

		for i := 0; i < count; i++ {
			p.replicate(<-p.orderRepEventChan)
			time.Sleep(INTERNAL * time.Millisecond)
		}
	}
}

func (p *ReplicateProcessor) replicate(event *e.OrderRepEvent) { //a. b: 备份
	//获取备份策略生成备份事件，马上备份。
	state, err := p.stateMachine.GetOrder(event.OrderId, dict.TASK_REP_FAIL)
	if err != nil {
		utils.Log(utils.WARN, "ReplicateProcessor  stateMachine.GetRepFailOrder ", err.Error(), event)
		return
	}

	strategy, err := utils.GetReplicationStrategy(event.OrderId)
	if err != nil {
		p.addRepEventToCache(event)
		utils.Log(utils.WARN, "ReplicateProcessor  GetReplicationStrategy ", err.Error(), event)
		return
	}

	repRequest, err := p.generateRepRequest(event.OrderId, state, strategy)
	if err != nil {
		utils.Log(utils.WARN, "ReplicateProcessor replicate generateRepRequest ", err.Error(), event)
		return
	}

	rsp, err := utils.Replicate(repRequest)
	if err != nil {
		p.addRepEventToCache(event)
		return
	}

	p.updateOrderState(rsp, state)

	if err := p.stateMachine.UpdateOrderStateInfo(repRequest.OrderId, state); err != nil {
		utils.Log(utils.WARN, "ReplicateProcessor replicate", err.Error(), state)
	}

	if state.Status != dict.TASK_BEGIN_REP { //还有区域没备份成功
		p.addRepEventToCache(event)
		return
	}

	utils.Log(utils.WARN, "ReplicateProcessor replicate success", "", rsp)
	//全部备份成功了,重新执行查询任务.
	p.orderChan <- repRequest.OrderId
}

func (p *ReplicateProcessor) updateOrderState(rsp *param.ReplicationResponse, state *dict.OrderStateInfo) {
	for _, task := range rsp.Tasks {
		stateTask, ok := state.Tasks[task.Fid]
		if !ok {
			continue
		}

		for region, status := range task.RegionStatus {
			if rep, ok := stateTask.Reps[region]; ok {
				if status != param.SUCCESS {
					rep.Status = dict.TASK_REP_FAIL
					stateTask.Status = dict.TASK_REP_FAIL //任务备份失败.
				} else {
					rep.Status = dict.TASK_BEGIN_REP
				}
			}
		}

		//遍历任务状态，如果有一个region 没备份成功，则备份失败。
		stateTask.Status = dict.TASK_BEGIN_REP
		for _, rep := range stateTask.Reps {
			if rep.Status == dict.TASK_REP_FAIL {
				stateTask.Status = dict.TASK_REP_FAIL
				break
			}
		}
	}

	state.Status = dict.TASK_BEGIN_REP //任务开始备份.
	for _, task := range state.Tasks {
		if task.Status != dict.TASK_BEGIN_REP { //还是备份失败.
			state.Status = dict.TASK_REP_FAIL
			break
		}
	}
}

func (p *ReplicateProcessor) addRepEventToCache(event *e.OrderRepEvent) {
	event.Count++
	if event.Count < dict.REP_COUNT {

		size := len(p.orderRepEventChan)
		if size >= p.repChanSize-1 {
			utils.Log(utils.ERROR, "ReplicateProcessor addRepEventToCache", fmt.Sprintf(" add  orderRepEventChan size = %v  have fill ", size), event)
		}

		p.orderRepEventChan <- event
	} else {
		utils.Log(utils.ERROR, "ReplicateProcessor addRepEventToCache", fmt.Sprintf("replicate count:%d more than %d", event.Count, dict.REP_COUNT), event)
	}
}

func (p *ReplicateProcessor) initRepEvent(orderIds []string) {
	for _, orderId := range orderIds {
		p.addRepEventToCache(&e.OrderRepEvent{Count: 0, OrderId: orderId})
	}
}

//to do get replicate strategy, 3 => 4, 3 不用备份， 4.备份失败, status == success. delete, 成功了，没有成功, 执行3备份, 4.success.
//get strategy
func (p *ReplicateProcessor) generateRepRequest(orderId string, order *dict.OrderStateInfo, strategy *param.StrategyInfo) (*param.ReplicationRequest, error) {
	request := &param.ReplicationRequest{
		OrderId: order.OrderId,
		Tasks:   make(map[string]*dict.Task, 10),
	}

	for _, val := range strategy.Tasks { //指定fid 的备份策略，strategy不用全部扫描。
		if task, ok := order.Tasks[val.Fid]; ok { //copy
			reqTaskReps := make(map[string]*dict.Rep, 3)
			for region, rep := range task.Reps {
				if rep.Status != dict.TASK_REP_FAIL {
					continue
				}

				//根据策略重新设置任务的分区备份数.
				if stRep, ok := val.Reps[region]; ok {
					//rep = stRep,直接赋值会改变订单状态。
					rep.MinRep = stRep.MinRep
					rep.MaxRep = stRep.MaxRep
				}

				reqRep := *rep
				reqTaskReps[region] = &reqRep
			}

			if len(reqTaskReps) > 0 {
				requsetTask := *task
				request.Tasks[task.Fid] = &requsetTask
			}
		}
	}
	return request, nil
}

func (p *ReplicateProcessor) processUploadFinishOrder(orders []string) error {
	for _, orderId := range orders {
		event := p.generateOrderUploadFinishEvent(orderId, param.SUCCESS)

		if err := p.stateMachine.Send(orderId, event); err != nil { //去获取策略，重新备份.
			utils.Log(utils.WARN, "ReplicateProcessor processUploadFinishOrder  p.stateMachine.Send fail ", err.Error(), "")
			continue
		}
		status := <-event.Ret
		if status != dict.SUCESS { //.添加到消息队列,重新备份。
			orderStatus, err := p.stateMachine.GetOrderStatus(orderId)
			if err == nil && orderStatus == dict.TASK_REP_FAIL { //重新备份。
				p.Add(orderId)
			}
		}
	}
	return nil
}

func (p *ReplicateProcessor) generateOrderUploadFinishEvent(orderId string, status int) *e.Event {
	return &e.Event{
		Type:    e.TASK_UPLOAD_FINISH,
		OrderId: orderId,
		Ret:     make(chan int),
		Data: &e.OrderUploadFinishEvent{
			OrderId: orderId,
			Status:  status,
		},
	}
}
