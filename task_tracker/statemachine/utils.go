package statemachine

import (
	"controller/task_tracker/dict"
	"controller/task_tracker/param"
	"errors"
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

func SetOrderState(state *dict.OrderStateInfo, status int) {
	for _, task := range state.Tasks {
		for _, rep := range task.Reps {
			rep.Status = status
		}
		task.Status = status
	}
	state.Status = status
}
