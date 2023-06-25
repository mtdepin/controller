package param

import "controller/strategy/dict"

const (
	SUCCESS = 1
	FAIL    = 0
)

type GetStrategyRequset struct {
	OrderId string `json:"order_id"`
}

type GetDeleteStrategyRequset struct {
	OrderId string `json:"order_id"`
}

type CreateStrategyRequest struct {
	RequestId string `json:"request_id"`
	OrderId   string `json:"order_id"`
}

type CreateStrategyResponse struct {
	Status int `json:"status"`
}

type GetDeleteStrategyResponse struct {
	Status  int    `json:"status"`
	OrderId string `json:"order_id"`
}

type StrategyInfo struct {
	Tasks []*dict.Task `json:"tasks"`
}

type GetStrategyResponse struct {
	Status   int           `json:"status"`
	OrderId  string        `json:"order_id"`
	Strategy *StrategyInfo `json:"strategy"`
}
