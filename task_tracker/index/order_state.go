package index

import (
	"controller/pkg/logger"
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	"errors"
	"fmt"
	"sync"
	"time"
)

type OrderStateIndex struct {
	orderStateMap map[string]*dict.OrderStateInfo
	store         *database.OrderState
	rwLock        *sync.RWMutex
}

func (p *OrderStateIndex) Init(db *database.DataBase) {
	p.rwLock = new(sync.RWMutex)
	p.store = new(database.OrderState)
	p.store.Init(db)

	states, err := p.store.Load()
	if err != nil {
		panic(fmt.Sprintf("load orderStateInfo fail, %v", err.Error()))
	}

	p.createIndex(*states)
}

func (p *OrderStateIndex) GetState(orderId string) (*dict.OrderStateInfo, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	if state, ok := p.orderStateMap[orderId]; ok {
		return p.copy(state), nil
	}
	return &dict.OrderStateInfo{}, errors.New("orderId not exist")
}

func (p *OrderStateIndex) Update(orderId string, state *dict.OrderStateInfo) error {
	p.rwLock.Lock()
	p.orderStateMap[orderId] = state
	p.rwLock.Unlock()

	return p.store.Update(state)
}

func (p *OrderStateIndex) SetTaskUploadStatus(orderId, fid, cid, region, origins string, status int) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	state, ok := p.orderStateMap[orderId]
	if !ok {
		return errors.New(fmt.Sprintf("OrderStateIndex SetTaskUploadStatus cid :%v err, fid: %v, orderId: %v, not find ", cid, fid, orderId))
	}

	if task, ok := state.Tasks[fid]; ok {
		task.Cid = cid
		task.Region = region
		task.Status = status
		task.Origins = origins
		state.UpdateTime = time.Now().UnixMilli()
		//return p.store.Update(state)
	} else {
		return errors.New(fmt.Sprintf("OrderStateIndex SetTaskUploadStatus cid :%v err, fid: %v not find , orderId: %v ", cid, fid, orderId))
	}

	bFlag := true
	for _, task := range state.Tasks {
		if task.Status != status {
			bFlag = false
			break
		}
	}
	if bFlag { //全部任务状态为 status ,则订单设置为status.
		state.Status = status
	}

	return p.store.Update(state)
}

func (p *OrderStateIndex) SetTaskStatus(orderId, fid, region string, status int) error {
	p.rwLock.Lock()
	state, ok := p.orderStateMap[orderId]
	if !ok {
		p.rwLock.Unlock()
		return errors.New(fmt.Sprintf("OrderStateIndex SetTaskStatus err, orderId: %v, not find ", orderId))
	}

	task, ok1 := state.Tasks[fid]
	if !ok1 {
		p.rwLock.Unlock()
		return errors.New(fmt.Sprintf("OrderStateIndex SetTaskStatus err, orderId: %v,  fid: %v not find ", orderId, fid))
	}

	rep, ok2 := task.Reps[region]
	if !ok2 {
		p.rwLock.Unlock()
		return errors.New(fmt.Sprintf("OrderStateIndex SetTaskStatus err, orderId: %v,  fid: %v, region: %v not find ", orderId, fid, region))
	}
	rep.Status = status
	state.UpdateTime = time.Now().UnixMilli()

	//如果所有任务每个区域都执行成功了, 设置任务为成功状态.
	bOrderFinish := true
	for _, task := range state.Tasks {

		bTaskFinish := true
		for _, rep := range task.Reps {
			if rep.Status != status { //任意区域没备份成功，都返回。
				bTaskFinish = false
				bOrderFinish = false
				break
			}
		}

		if bTaskFinish == false {
			continue
		}
		task.Status = status //任务所有区域都执行成功了。
	}

	if bOrderFinish {
		state.Status = status
	}

	p.rwLock.Unlock()
	return p.store.Update(state)
}

func (p *OrderStateIndex) SetStatus(orderId string, status int) error {
	p.rwLock.Lock()
	state, ok := p.orderStateMap[orderId]
	if !ok {
		p.rwLock.Unlock()
		return errors.New(fmt.Sprintf("OrderStateIndex SetChargeStatus err, orderId: %v, not find ", orderId))
	}

	state.Status = status
	state.UpdateTime = time.Now().UnixMilli()

	p.rwLock.Unlock()
	return p.store.Update(state)
}

func (p *OrderStateIndex) createIndex(states []dict.OrderStateInfo) {
	p.orderStateMap = make(map[string]*dict.OrderStateInfo, len(states)+10)
	for i, _ := range states {
		p.orderStateMap[states[i].OrderId] = &states[i]
	}
}

func (p *OrderStateIndex) copy(src *dict.OrderStateInfo) *dict.OrderStateInfo {
	if src == nil {
		return nil
	}
	dest := new(dict.OrderStateInfo)
	dest.OrderId = src.OrderId
	dest.OrderType = src.OrderType
	dest.Status = src.Status
	dest.CreateTime = src.CreateTime
	dest.UpdateTime = src.UpdateTime

	dest.Tasks = make(map[string]*dict.Task, len(src.Tasks))
	for key, task := range src.Tasks {
		destTask := &dict.Task{
			Fid:     task.Fid,
			Cid:     task.Cid,
			Region:  task.Region,
			Origins: task.Origins,
			Status:  task.Status,
		}
		destTask.Reps = make(map[string]*dict.Rep, len(task.Reps))
		for region, rep := range task.Reps {
			destTask.Reps[region] = &dict.Rep{
				Region:     rep.Region,
				VirtualRep: rep.VirtualRep,
				RealRep:    rep.RealRep,
				MinRep:     rep.MinRep,
				MaxRep:     rep.MaxRep,
				Expire:     rep.Expire,
				Encryption: rep.Encryption,
				Status:     rep.Status,
			}
		}
		dest.Tasks[key] = destTask
	}

	return dest
}

func (p *OrderStateIndex) GetOrderStatus(orderId string) (int, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	state, ok := p.orderStateMap[orderId]
	if !ok {
		return 0, errors.New(fmt.Sprintf("OrderStateIndex GetOrderStatus err, orderId: %v, not find ", orderId))
	}

	return state.Status, nil
}

func (p *OrderStateIndex) DeleleOrder(orderId string) error {
	p.rwLock.Lock()
	defer p.rwLock.Unlock()
	if _, ok := p.orderStateMap[orderId]; !ok {
		return errors.New(fmt.Sprintf("OrderStateIndex DeleleOrder err, orderId: %v, not find ", orderId))
	}

	delete(p.orderStateMap, orderId)
	return nil
}

func (p *OrderStateIndex) UpdateOrderRepInfo(orderId string, tasks map[string]*dict.TaskRepInfo) error {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	orderStateInfo, ok := p.orderStateMap[orderId]
	if !ok {
		return errors.New(fmt.Sprintf("OrderStateIndex UpdateOrderRepInfo fail, orderid: %v no find", orderId))
	}

	bUpdate := false
	for _, task := range tasks {
		orderTask, ok := orderStateInfo.Tasks[task.Fid]
		if !ok {
			logger.Warnf("OrderStateIndex UpdateOrderRepInfo fail, orderid: %v ,  the fid: %v not find", orderId, task.Fid)
			continue
		}

		for _, regionRep := range task.Regions {
			if regionRep.Status != dict.SUCESS {
				continue
			}

			rep, ok := orderTask.Reps[regionRep.Region]
			if !ok {
				logger.Warnf("OrderStateIndex UpdateOrderRepInfo fail, orderid: %v ,   fid: %v  the region: %v not find , ", orderId, task.Fid, regionRep.Region)
				continue
			}
			//存在
			if rep.RealRep != 0 {
				if rep.RealRep == regionRep.CurRep {
					rep.Status = dict.TASK_REP_SUC
					regionRep.CheckStatus = dict.SUCESS
					bUpdate = true
				}
			} else { //rep.RealRep == 0
				if regionRep.CurRep >= rep.MinRep {
					rep.RealRep = regionRep.CurRep
					rep.Status = dict.TASK_REP_SUC
					regionRep.CheckStatus = dict.SUCESS
					bUpdate = true
				}
			}
		}

		//判断任务的所有region是否备份成功.
		taskRepSuc := true
		for _, rep := range orderTask.Reps {
			if rep.Status != dict.TASK_REP_SUC {
				taskRepSuc = false
				break
			}
		}
		if taskRepSuc {
			orderTask.Status = dict.TASK_REP_SUC
		}
	}

	if !bUpdate { //如果没有更新直接返回
		return nil
	}

	//判断订单的所有任务是否备份成功
	orderRepSuc := true
	for _, task := range orderStateInfo.Tasks {
		if task.Status != dict.TASK_REP_SUC {
			orderRepSuc = false
			break
		}
	}

	if orderRepSuc {
		orderStateInfo.Status = dict.TASK_REP_SUC
	}

	// to do save
	return p.store.Update(orderStateInfo)
}

func (p *OrderStateIndex) GetAllUploadFinishOrderInfo() []*dict.UploadFinishOrder {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	orders := make([]*dict.UploadFinishOrder, 0, len(p.orderStateMap))
	for _, state := range p.orderStateMap {
		if state.Status == dict.TASK_UPLOAD_SUC {
			order := p.getFinishOrderInfo(state)
			orders = append(orders, order)
		}
	}
	return orders
}

func (p *OrderStateIndex) GetUploadFinishOrderInfo(orderId string) (*dict.UploadFinishOrder, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	if state, ok := p.orderStateMap[orderId]; ok {
		if state.Status == dict.TASK_UPLOAD_SUC {
			return p.getFinishOrderInfo(state), nil
		} else {
			return nil, errors.New(fmt.Sprintf("OrderStateIndex, GetUploadFinishOrderInfo, fail , orderId: %v  state: %d not upload finish ", orderId, state.Status))
		}
	}

	return nil, errors.New(fmt.Sprintf("OrderStateIndex, GetUploadFinishOrderInfo, fail , orderId: %v not find", orderId))
}

func (p *OrderStateIndex) getFinishOrderInfo(state *dict.OrderStateInfo) *dict.UploadFinishOrder {
	order := &dict.UploadFinishOrder{
		OrderId: state.OrderId,
		Tasks:   make([]*dict.RepTask, 0, len(state.Tasks)),
	}

	for fid, val := range state.Tasks {
		taskRequest := &dict.RepTask{
			Fid:     fid,
			Cid:     val.Cid,
			Regions: make([]string, 0, len(val.Reps)),
		}

		for region, _ := range val.Reps {
			taskRequest.Regions = append(taskRequest.Regions, region)
		}

		order.Tasks = append(order.Tasks, taskRequest)
	}

	return order
}

func (p *OrderStateIndex) GetAllOrders(status int) []*dict.OrderStateInfo {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	orders := make([]*dict.OrderStateInfo, 0, len(p.orderStateMap))
	for _, state := range p.orderStateMap {
		if state.Status == status {
			order := p.copy(state)
			orders = append(orders, order)
		}
	}
	return orders
}

func (p *OrderStateIndex) GetOrder(orderId string, status int) (*dict.OrderStateInfo, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	if state, ok := p.orderStateMap[orderId]; ok {
		if state.Status == status {
			return p.copy(state), nil
		} else {
			return nil, errors.New(fmt.Sprintf("OrderStateIndex, GetOrder, fail , orderId: %v  status: %d != TASK_REP_FAIL ", orderId, state.Status))
		}
	}

	return nil, errors.New(fmt.Sprintf("OrderStateIndex, GetOrder, fail , orderId: %v not find", orderId))
}

func (p *OrderStateIndex) GetFidStatus(orderId, fid string) (int, error) {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	state, ok := p.orderStateMap[orderId]

	if !ok {
		return 0, errors.New(fmt.Sprintf("OrderStateIndex GetFidStatus err, orderId: %v, not find ", orderId))
	}

	task, ok1 := state.Tasks[fid]
	if !ok1 {
		return 0, errors.New(fmt.Sprintf("OrderStateIndex GetFidStatus err,orderId: %v, fid: %v, not find ", orderId, fid))
	}

	return task.Status, nil
}
