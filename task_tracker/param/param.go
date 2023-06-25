package param

import "controller/task_tracker/dict"

const (
	FAIL    = 0
	SUCCESS = 1
	PROCEED = 2 //处理中
)

const (
	HaveDevice = 1
	NoDevice   = 2
)

const (
	UPLOAD   = 1
	DOWNLOAD = 2
)

const (
	CHANAL_SIZE        = 1024 * 8
	CHARGE_CHANAL_SIZE = 1024 * 8
	ORDER_CACHE_SIZE   = 1024 * 5
)

type CreateTaskRequest struct {
	RequestId string `json:"request_id"`
	Type      int    `json:"task_type"`
}

type OrderTaskResponse struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"`
}

type CreateStrategyRequest struct {
	RequestId string `json:"request_id"`
	OrderId   string `json:"order_id"`
}

type CreateStrategyResponse struct {
	Status     int    `json:"status"`
	StrategyId string `json:"strategy_id"`
}

type UploadTaskResponse struct {
	Status  int    `json:"status"`
	OrderId string `json:"order_id"`
}

type CheckBalanceRequest struct {
	UserId string `json:"user_id"`
}

type CheckBalanceResponse struct {
	Status int  `json:"status"`
	Enough bool `json:"enough"`
}

type UploadTask struct {
	Fid     string `json:"fid"`
	Cid     string `json:"cid"`
	Region  string `json:"region"`
	Origins string `json:"origins"`
	Status  int    `json:"status"`
}

type UploadFinishRequest struct {
	OrderId string        `json:"order_id"`
	Tasks   []*UploadTask `json:"tasks"`
	Status  int           `json:"status"`
}

type UploadFinishResponse struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"`
}

type CallbackUploadRequest struct {
	OrderId string `json:"order_id"`
	Fid     string `json:"fid"`
	Cid     string `json:"cid"`
	Region  string `json:"region"`
	Origins string `json:"origins"`
	Status  int    `json:"status"`
}

type CallbackUploadResponse struct {
	Status int `json:"status"` //1:成功， 0:失败
}

type GetStrategyRequset struct {
	OrderId string `json:"order_id"`
}

type TaskResponse struct {
	Fid          string         `json:"fid"` //文件hash值
	Cid          string         `json:"cid"`
	RegionStatus map[string]int `json:"region_status"` //key region, value status
}

type StrategyInfo struct {
	Tasks []*dict.Task `json:"tasks"`
}

type GetStrategyResponse struct {
	Status   int           `json:"status"`
	OrderId  string        `json:"order_id"`
	Strategy *StrategyInfo `json:"strategy"`
}

type ReplicationRequest struct {
	OrderId string                `json:"order_id"`
	Origins string                `json:"origins"`
	Tasks   map[string]*dict.Task `json:"tasks"` //key fid
}

type ReplicationResponse struct {
	Status  int             `json:"status"`
	OrderId string          `json:"order_id"`
	Tasks   []*TaskResponse `json:"tasks"`
}

type DeleteOrderRequest struct {
	OrderId string                 `json:"order_id"`
	Tasks   map[string]*UploadTask `json:"tasks"` //key fid
}

type DeleteOrderResponse struct {
	Status  int           `json:"status"`
	OrderId string        `json:"order_id"`
	Tasks   []*UploadTask `json:"tasks"`
}

type CallbackRepRequest struct {
	OrderId string `json:"order_id"`
	Region  string `json:"region"`
	Fid     string `json:"fid"`
	Cid     string `json:"cid"`
	Status  int    `json:"status"`
}

type CallbackRepResponse struct {
	Status int `json:"status"`
}

type CallbackDelResponse struct {
	Status int `json:"status"`
}

type CallbackDeleteRequest struct {
	OrderId string `json:"order_id"`
	Fid     string `json:"fid"`
	Cid     string `json:"cid"`
	Region  string `json:"region"`
	Status  int    `json:"status"`
}

type ChargeRequest struct {
	OrderId   string       `json:"order_id"`
	OrderType int32        `json:"order_type"`
	Tasks     []*dict.Task `json:"tasks,omitempty"`
}

type ChargeResponse struct {
	Status int `json:"status"`
}

type CallbackChargeRequest struct {
	OrderId   string `json:"order_id"`
	OrderType int    `json:"order_type"`
	Status    int    `json:"status"`
}

type CallbackChargeResponse struct {
	Status int `json:"status"`
}

type GetOrderRepResponse struct {
	Status  int                          `json:"status"`
	OrderId string                       `json:"order_id"`
	Tasks   map[string]*dict.TaskRepInfo `json:"tasks"`
}

type DownloadFinishRequest struct {
	OrderId string         `json:"order_id"`
	Tasks   map[string]int `json:"tasks"` //key cid, value, status.
	Status  int            `json:"status"`
}

type DownloadFinishResponse struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"`
}

type DownloadTaskResponse struct {
	Status  int    `json:"status"`
	OrderId string `json:"order_id"`
}
