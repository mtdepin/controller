package event

import (
	"controller/task_tracker/dict"
	"controller/task_tracker/param"
)

const (
	CREATE_ORDER = iota
	CALLBACK_UPLOAD
	TASK_UPLOAD_FINISH
	CALLBACK_REP
	CALLBACK_CHARGE
	TASK_DOWNLOAD_FINISH
)

type RepInfo struct {
	Region     int //备份区域号
	VirtualRep int //虚拟备份数
	RealRep    int //实际备份数
	MinRep     int //最小备份数
	MaxRep     int //最大备份数
	Expire     int //过期时间
	Encryption int //是否加密
	Status     int //备份状态.
}

type Task struct {
	Fid    string             //文件hash值
	Cid    string             //文件cid
	RepMap map[string]RepInfo //key: 表示区域信息， value: 文件区域备份信息
	Status int                //文件状态
}

type CallbackUploadEvent struct {
	OrderId string `json:"order_id"`
	Fid     string `json:"fid"`
	Cid     string `json:"cid"`
	Region  string `json:"region"`
	Origins string `json:"origins"`
	Status  int    `json:"status"`
}

type CallbackRepEvent struct {
	OrderId string `json:"order_id"`
	Fid     string `json:"fid"`
	Cid     string `json:"cid"`
	Region  string `json:"region"`
	Status  int    `json:"status"`
}

type ReplicateEvent struct {
	OrderId        string            `json:"order_id"`
	Cid            string            `json:"cid"`
	Meta           map[string]string `json:"dict,omitempty"`
	Replication    int               `json:"replication"`
	ReplicationMin int               `json:"replication_min"`
	ReplicationMax int               `json:"replication_max"`
	Allocations    []string          `json:"allocations,omitempty"`
	Expire         uint64            `json:"expire"`
}

type DeleteEvent struct {
	OrderId string `json:"order_id"`
	Cid     string `json:"cid"`
}

type CallbackDeleteEvent struct {
	OrderId string `json:"order_id"`
	Fid     string `json:"fid"`
	Cid     string `json:"cid"`
	Region  string `json:"region"`
	Status  int    `json:"status"`
}

type CallbackChargeEvent struct {
	OrderId   string `json:"order_id"`
	OrderType int    `json:"order_type"`
	Status    int    `json:"status"`
}

type OrderUploadFinishEvent struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"` //1.上传成功， 0.上传失败
}

type OrderDownloadFinishEvent struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"` //1.上传成功， 0.上传失败
}

type CreateOrderEvent struct {
	RequestId string   `json:"request_id"`
	OrderType int      `json:"order_type"`
	Fids      []string `json:"fids"`
	Cids      []string `json:"cids"`
}

type Event struct {
	Type    int //事件类型
	OrderId string
	Ret     chan int
	Data    interface{}
}

type OrderRepCheckEvent struct {
	Count     int   //请求执行次数
	BeginTime int64 //开始查询时间
	Request   *dict.UploadFinishOrder
}

type OrderRepEvent struct {
	Count   int //请求执行次数
	Request *param.ReplicationRequest
}

type OrderDeleteEvent struct {
	Count   int //请求执行次数
	Request *param.DeleteOrderRequest
}

type OrderChargeEvent struct {
	Count   int //请求执行次数
	Request *param.ChargeRequest
}
