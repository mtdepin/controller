package statemachine

import (
	"controller/task_tracker/dict"
	"controller/task_tracker/param"
	"errors"
	"time"
)

func generateChargeRequest(orderId string, state *dict.OrderStateInfo) (*param.ChargeRequest, error) {
	if state == nil {
		return nil, errors.New("generateChargeRequest OrderStateInfo is nil")
	}

	request := &param.ChargeRequest{OrderId: orderId, OrderType: state.OrderType, Tasks: make([]*dict.Task, 0, len(state.Tasks))}
	for _, task := range state.Tasks {
		request.Tasks = append(request.Tasks, task)
	}

	return request, nil
}

func setOrderReplicateInfo(strategy *param.StrategyInfo, state *dict.OrderStateInfo) error {
	if strategy == nil || state == nil {
		return errors.New("param error strategy or state is nil")
	}

	for _, val := range strategy.Tasks {
		if task, ok := state.Tasks[val.Fid]; ok {
			task.Reps = make(map[string]*dict.Rep, len(val.Reps))
			//if fid replicate success,则不必在备份
			for _, rep := range val.Reps {
				task.Reps[rep.Region] = &dict.Rep{
					Region:     rep.Region,
					VirtualRep: rep.VirtualRep,
					RealRep:    rep.RealRep,
					MinRep:     rep.MinRep,
					MaxRep:     rep.MaxRep,
					Expire:     rep.Expire,
					Encryption: rep.Encryption,
					Status:     dict.TASK_REP_FAIL, //默认备份失败
				}
			}
			task.Status = dict.TASK_REP_FAIL
		}
	}

	state.Status = dict.TASK_REP_FAIL
	state.UpdateTime = time.Now().UnixMilli()

	return nil
}

func setOrderRspState(rsp *param.ReplicationResponse, state *dict.OrderStateInfo) {
	for _, task := range rsp.Tasks {
		stateTask, ok := state.Tasks[task.Fid]
		if !ok {
			continue
		}

		for region, status := range task.RegionStatus {
			if rep, ok := stateTask.Reps[region]; ok { //可能只响应一部分.
				if status == param.SUCCESS {
					rep.Status = dict.TASK_BEGIN_REP
				}
			}
		}

		//一个集群失败，则fid任务状态为备份失败.
		stateTask.Status = dict.TASK_BEGIN_REP
		for _, rep := range stateTask.Reps {
			if rep.Status == dict.TASK_REP_FAIL {
				stateTask.Status = dict.TASK_REP_FAIL
				break
			}
		}
	}

	state.Status = dict.TASK_BEGIN_REP
	for _, task := range state.Tasks { //一个任务备份失败，则集群备份失败
		if task.Status != dict.TASK_BEGIN_REP {
			state.Status = dict.TASK_REP_FAIL
			break
		}
	}
}
