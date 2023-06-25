package statemachine

import (
	"controller/task_tracker/dict"
	e "controller/task_tracker/event"
	"controller/task_tracker/index"
	"controller/task_tracker/param"
	"controller/task_tracker/utils"
	"errors"
	"time"
)

type CallbackUpload struct {
	orderIndex      *index.OrderIndex
	orderStateIndex *index.OrderStateIndex
}

func (p *CallbackUpload) Init(orderIndex *index.OrderIndex, orderStateIndex *index.OrderStateIndex) {
	p.orderIndex = orderIndex
	p.orderStateIndex = orderStateIndex
}

func (p *CallbackUpload) HandleUploadEvent(event *e.Event) error {
	callbackUploadEvent := event.Data.(*e.CallbackUploadEvent)

	if callbackUploadEvent.Status == param.SUCCESS {
		return p.handleUploadSucEvent(callbackUploadEvent)
	}
	return p.handleUploadFailEvent(callbackUploadEvent)
}

func (p *CallbackUpload) handleUploadSucEvent(event *e.CallbackUploadEvent) error {
	if err := p.orderStateIndex.SetTaskUploadStatus(event.OrderId, event.Fid, event.Cid, event.Region, event.Origins, dict.TASK_UPLOAD_SUC); err != nil {
		return err
	}

	status, err := p.orderStateIndex.GetOrderStatus(event.OrderId)
	if err != nil || status != dict.TASK_UPLOAD_SUC {
		return err
	}

	//文件上传成功，开始备份.
	return p.Replicate(event.OrderId)
}

func (p *CallbackUpload) handleUploadFailEvent(event *e.CallbackUploadEvent) error {
	return p.orderStateIndex.SetTaskUploadStatus(event.OrderId, event.Fid, event.Cid, event.Region, event.Origins, dict.TASK_UPLOAD_FAIL)
}

//备份文件
func (p *CallbackUpload) Replicate(orderId string) error {
	strategy, err := utils.GetReplicationStrategy(orderId)
	if err != nil {
		return err
	}

	state, err1 := p.orderStateIndex.GetState(orderId)
	if err1 != nil {
		return err1
	}

	if err := p.createOrderUploadFinishState(strategy, state); err != nil {
		return err
	}

	if err := p.orderStateIndex.Update(orderId, state); err != nil {
		return err
	}

	if err := p.orderIndex.UpdateStatus(orderId, dict.TASK_UPLOAD_SUC); err != nil {
		return err
	}

	if rsp, err := utils.Replicate(&param.ReplicationRequest{OrderId: orderId, Tasks: state.Tasks}); err != nil {
		//订单所有任务都设置成失败.
		SetOrderState(state, dict.TASK_REP_FAIL)
		utils.Log(utils.WARN, "handleUploadFinishSucEvent replicate ", err.Error(), &param.ReplicationRequest{OrderId: orderId, Tasks: state.Tasks})
	} else {
		p.setOrderState(rsp, state)
	}

	return p.orderStateIndex.Update(orderId, state)
}

func (p *CallbackUpload) setTaskCid(strategy *param.StrategyInfo, state *dict.OrderStateInfo) {
	for _, val := range strategy.Tasks {
		if task, ok := state.Tasks[val.Fid]; ok {
			val.Cid = task.Cid
		}
	}
}

func (p *CallbackUpload) createOrderUploadFinishState(strategy *param.StrategyInfo, state *dict.OrderStateInfo) error {
	if strategy == nil || state == nil {
		return errors.New("param error strategy or state is nil")
	}

	for _, val := range strategy.Tasks {
		if task, ok := state.Tasks[val.Fid]; ok {
			task.Reps = make(map[string]*dict.Rep, len(val.Reps))
			for _, rep := range val.Reps {
				task.Reps[rep.Region] = &dict.Rep{
					Region:     rep.Region,
					VirtualRep: rep.VirtualRep,
					RealRep:    rep.RealRep,
					MinRep:     rep.MinRep,
					MaxRep:     rep.MaxRep,
					Expire:     rep.Expire,
					Encryption: rep.Encryption,
					Status:     dict.TASK_UPLOAD_SUC,
				}
			}
			task.Status = dict.TASK_UPLOAD_SUC
		}
	}

	state.Status = dict.TASK_UPLOAD_SUC //订单上传完成
	state.UpdateTime = time.Now().UnixMilli()

	return nil
}

func (p *CallbackUpload) setOrderState(rsp *param.ReplicationResponse, state *dict.OrderStateInfo) {
	for _, task := range rsp.Tasks {
		stateTask, ok := state.Tasks[task.Fid]
		if !ok {
			continue
		}

		for region, status := range task.RegionStatus {
			if status != param.SUCCESS {
				if rep, ok := stateTask.Reps[region]; ok {
					rep.Status = dict.TASK_REP_FAIL
					state.Status = dict.TASK_REP_FAIL //订单中有一个区域没备份成功，则订单备份失败.
				}
			}
		}
	}
}
