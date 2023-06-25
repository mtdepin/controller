package param

import (
	"controller/task_tracker/dict"
	"golang.org/x/net/context"
)

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

type CreateStrategyRequest struct {
	RequestId string  `json:"request_id"`
	OrderId   string  `json:"order_id"`
	Region    string  `json:"region"`
	Ext       *Extend `json:"ext,omitempty"`
}

type CreateStrategyResponse struct {
	Status     int    `json:"status"`
	StrategyId string `json:"strategy_id"`
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
	Ext     *Extend       `json:"ext,omitempty"`
}

type UploadFinishResponse struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"`
}

type CallbackUploadRequest struct {
	OrderId string  `json:"order_id"`
	Fid     string  `json:"fid"`
	Cid     string  `json:"cid"`
	Region  string  `json:"region"`
	Origins string  `json:"origins"`
	Status  int     `json:"status"`
	Ext     *Extend `json:"ext,omitempty"`
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
	Ext     *Extend                `json:"ext,omitempty"`
}

type DeleteOrderResponse struct {
	Status  int           `json:"status"`
	OrderId string        `json:"order_id"`
	Tasks   []*UploadTask `json:"tasks"`
}

type CallbackRepRequest struct {
	OrderId string  `json:"order_id"`
	Region  string  `json:"region"`
	Fid     string  `json:"fid"`
	Cid     string  `json:"cid"`
	Status  int     `json:"status"`
	Ext     *Extend `json:"ext,omitempty"`
}

type CallbackRepResponse struct {
	Status int `json:"status"`
}

type CallbackDelResponse struct {
	Status int `json:"status"`
}

type CallbackDeleteRequest struct {
	OrderId string  `json:"order_id"`
	Fid     string  `json:"fid"`
	Cid     string  `json:"cid"`
	Region  string  `json:"region"`
	Status  int     `json:"status"`
	Ext     *Extend `json:"ext,omitempty"`
}

type ChargeRequest struct {
	OrderId   string       `json:"order_id"`
	OrderType int32        `json:"order_type"`
	Tasks     []*dict.Task `json:"tasks,omitempty"`
	Ext       *Extend      `json:"ext,omitempty"`
}

type ChargeResponse struct {
	Status int `json:"status"`
}

type CallbackChargeRequest struct {
	OrderId   string  `json:"order_id"`
	OrderType int     `json:"order_type"`
	Status    int     `json:"status"`
	Ext       *Extend `json:"ext,omitempty"`
}

type CallbackChargeResponse struct {
	Status int `json:"status"`
}

type GetOrderRepResponse struct {
	Status  int                          `json:"status"`
	OrderId string                       `json:"order_id"`
	Tasks   map[string]*dict.TaskRepInfo `json:"tasks"` //可以采用动态生成的方式， 好处减少内存使用，不好浪费cpu.
}

type DownloadFinishRequest struct {
	OrderId string         `json:"order_id"`
	Tasks   map[string]int `json:"tasks"` //key cid, value, status.
	Status  int            `json:"status"`
	Ext     *Extend        `json:"ext,omitempty"`
}

type DownloadFinishResponse struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"`
}

type DownloadTaskResponse struct {
	Status  int    `json:"status"`
	OrderId string `json:"order_id"`
}

type DeleteFidRequest struct {
	RequestId string          `json:"request_id"`
	UserId    string          `json:"user_id"`
	OrderId   string          `json:"order_id"`
	Fids      map[string]bool `json:"fids"`
	Ext       *Extend         `json:"ext,omitempty"`
}

type GetFidDeleteStrategyRequest struct {
	OrderId string          `json:"order_id"`
	Fids    map[string]bool `json:"fids"`
	Ext     *Extend         `json:"ext,omitempty"`
}

type DeleteFidResponse struct {
	OrderId string         `json:"order_id"`
	Fids    map[string]int `json:"fids"` //每个fid的删除状态
	Status  int            `json:"status"`
}

//调度请求
type DeleteOrderFidRequest struct {
	OrderId string                   `json:"order_id"`
	Tasks   map[string][]*UploadTask `json:"tasks"` //key fid
	Ext     *Extend                  `json:"ext,omitempty"`
}

//调度响应.
type DeleteOrderFidResponse struct {
	Status  int                       `json:"status"`
	OrderId string                    `json:"order_id"`
	Tasks   map[string]*[]*UploadTask `json:"tasks"` //key fid
}

type CallbackDownloadRequest struct {
	Cid     string  `json:"cid"`
	Region  string  `json:"region"`
	Origins string  `json:"origins"`
	Status  int     `json:"status"`
	Ext     string  `json:"ext"`
	Extend  *Extend `json:"extend,omitempty"`
}

type CallbackDownloadResponse struct {
	Status int `json:"status"` //1:成功， 0:失败
}

type FidPageQueryRequest struct {
	PageField string  `json:"page_field"` //分页查询字段， 本次默认按create_time 降序查询.
	LastValue int64   `json:"last_value"` //上一页最后一个字段值
	Sort      int     `json:"sort"`       //排序方式，1.升序，2.降序
	PageSize  int     `json:"page_size"`
	Ext       *Extend `json:"ext,omitempty"`
}

type FidPageQueryResponse struct {
	FidInfos []*FidInfo `json:"fid_infos"`
	Status   int        `json:"status"`
}

type FidInfo struct {
	Fid        string `json:"fid"`
	Cid        string `json:"cid"`
	Origins    string `json:"origins"` //文件上传节点
	Status     int    `json:"status"`  //状态
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type OrderPageQueryRequest struct {
	OrderType int32   `json:"order_type"` //1.上传订单, 2.下载订单
	PageField string  `json:"page_field"` //分页查询字段， 本次默认按create_time 降序查询.
	LastValue int64   `json:"last_value"` //上一页最后一个字段值
	Sort      int     `json:"sort"`       //排序方式，1.升序，2.降序
	PageSize  int     `json:"page_size"`
	Ext       *Extend `json:"ext,omitempty"`
}

type OrderPageQueryResponse struct {
	Orders []*OrderInfo `json:"orders"`
	Status int          `json:"status"`
}

type OrderInfo struct {
	OrderId    string `json:"order_id"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type SearchOrderRequest struct {
	OrderIds  []string `json:"order_ids"`
	OrderType int      `json:"order_type"`
	Ext       *Extend  `json:"ext,omitempty"`
}

type SearchOrderResponse struct {
	Orders []*dict.OrderStateInfo `json:"orders"`
	Status int                    `json:"status"` //1:成功， 0:失败
}

type SearchFidRequest struct {
	Fids []string `json:"fids"`
	Ext  *Extend  `json:"ext,omitempty"`
}

type SearchFidResponse struct {
	FidInfos []*dict.FidInfo `json:"fidInfos"`
	Status   int             `json:"status"` //1:成功， 0:失败
}

type Extend struct {
	Ctx context.Context `json:"ctx,omitempty"`
}

type FidRepInfo struct {
	Fid     int `json:"fid"`
	RealRep int `json:"real_rep"`
	MinRep  int `json:"min_rep"`
	MaxRep  int `json:"max_rep"`
	Status  int `json:"status"`
}

type CreateSafeStrategyRequest struct {
	RequestId   string        `json:"request_id"`
	FidRepInfos []*FidRepInfo `json:"fid_rep_infos"`
	Region      string        `json:"region"`
}

type CreateSafeStrategyResponse struct {
	FidRepInfo *FidRepInfo `json:"fid_rep_info"`
}

type GetSafeStrategyResponse struct {
	UpdateOrderIds []string      `json:"update_order_ids"`
	OrderId        string        `json:"order_id"`
	Strategy       *StrategyInfo `json:"strategy"`
	Status         int           `json:"status"`
}

type GetDynamicStrategyRequest struct {
	Fid string `json:"fid"` //需要调整备份数的文件
}

//
type GetDynamicStrategyResponse struct {
	OrderId  string        `json:"order_id"`
	Strategy *StrategyInfo `json:"strategy"`
	Status   int           `json:"status"`
}

type GetBalanceStrategyRequest struct {
	Fid         string   `json:"fid"`          //需要迁移的文件
	SrcRegion   string   `json:"src_region"`   //文件当前备份区域
	DestRegions []string `json:"dest_regions"` //文件迁移的目标候选集群。
}

type GetBalanceStrategyResponse struct {
	UpdateOrderIds []string      `json:"update_order_ids"` //需要更新和fid关联的所有订单
	OrderId        string        `json:"order_id"`         //当前执行备份的订单号
	Strategy       *StrategyInfo `json:"strategy"`         //备份策略
	Status         int           `json:"status"`
}
