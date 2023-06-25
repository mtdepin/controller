package processor

import (
	"controller/pkg/logger"
	"controller/task_tracker/database"
	"controller/task_tracker/dict"
	"controller/task_tracker/param"
	"controller/task_tracker/statemachine"
)

type SearchProcessor struct {
	stateMachine     *statemachine.StateMachine
	orderStateDb     *database.OrderState
	fidReplicationDb *database.FidReplication
}

func (p *SearchProcessor) Init(machine *statemachine.StateMachine, db *database.DataBase) {
	p.stateMachine = machine
	p.orderStateDb = new(database.OrderState)
	p.orderStateDb.Init(db)
	p.fidReplicationDb = new(database.FidReplication)
	p.fidReplicationDb.Init(db)
}

func (p *SearchProcessor) GetOrderDetail(request *param.SearchOrderRequest) (interface{}, error) {
	orders := make([]*dict.OrderStateInfo, 0, len(request.OrderIds))

	for _, orderId := range request.OrderIds {
		if order, err := p.stateMachine.GetOrderStateInfo(orderId); err == nil {
			orders = append(orders, order)
		} else {
			if ret, err := p.orderStateDb.GetOrderStateByOrderId(orderId); err == nil {
				orders = append(orders, &(*ret)[0])
			} else {
				logger.Warnf("get orderId from db fail: %v", err.Error())
			}
		}
	}

	status := param.SUCCESS
	if len(orders) == 0 {
		status = param.FAIL
	}

	return &param.SearchOrderResponse{
		Orders: orders,
		Status: status,
	}, nil
}

func (p *SearchProcessor) GetPageOrders(request *param.OrderPageQueryRequest) (interface{}, error) {
	orderStates, err := p.orderStateDb.QueryOrderByPage(request.OrderType, request.LastValue, request.PageField, request.Sort, request.PageSize)
	if err != nil {
		logger.Warnf("SearchProcessor, GetPageOrders fail: %v", err.Error())
		return &param.OrderPageQueryResponse{
			Orders: make([]*param.OrderInfo, 0, 1),
			Status: param.FAIL,
		}, nil
	}

	orders := make([]*param.OrderInfo, 0, len(*orderStates))
	for _, order := range *orderStates {
		orders = append(orders,
			&param.OrderInfo{
				OrderId:    order.OrderId,
				Status:     order.Status,
				CreateTime: order.CreateTime,
				UpdateTime: order.UpdateTime,
			})
	}

	return &param.OrderPageQueryResponse{
		Orders: orders,
		Status: param.SUCCESS,
	}, nil
}

func (p *SearchProcessor) GetFidDetail(request *param.SearchFidRequest) (interface{}, error) {
	fidInfos := make([]*dict.FidInfo, 0, len(request.Fids))

	ret, err := p.fidReplicationDb.SearchFids(request.Fids)
	if err != nil {
		logger.Warnf("SearchProcessor, fidReplicationDb fail: %v", err.Error())
		return &param.SearchFidResponse{
			FidInfos: make([]*dict.FidInfo, 0, 1),
			Status:   param.FAIL,
		}, nil
	}

	for _, fidInfo := range *ret {
		fidInfos = append(fidInfos, &fidInfo)
	}

	return &param.SearchFidResponse{
		FidInfos: fidInfos,
		Status:   param.SUCCESS,
	}, nil
}

func (p *SearchProcessor) GetPageFids(request *param.FidPageQueryRequest) (interface{}, error) {
	ret, err := p.fidReplicationDb.QueryFidByPage(request.LastValue, request.PageField, request.Sort, request.PageSize)
	if err != nil {
		logger.Warnf("SearchProcessor, GetPageFids fail: %v", err.Error())
		return &param.FidPageQueryResponse{
			FidInfos: make([]*param.FidInfo, 0, 1),
			Status:   param.FAIL,
		}, nil
	}

	fidInfos := make([]*param.FidInfo, 0, len(*ret))
	for _, fidInfo := range *ret {
		status := dict.TASK_INIT
		if fidInfo.Cid != "" {
			status = dict.TASK_UPLOAD_SUC
		}

		fidInfos = append(fidInfos,
			&param.FidInfo{
				Fid:        fidInfo.Fid,
				Cid:        fidInfo.Cid,
				Origins:    fidInfo.Origins,
				Status:     status,
				CreateTime: fidInfo.CreateTime,
				UpdateTime: fidInfo.UpdateTime,
			})
	}

	return &param.FidPageQueryResponse{
		FidInfos: fidInfos,
		Status:   param.SUCCESS,
	}, nil
}
