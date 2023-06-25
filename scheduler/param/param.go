package param

import "golang.org/x/net/context"

const (
	FAIL    = 0
	SUCCESS = 1
)

type UploadTask struct {
	Fid     string `json:"fid"`
	Cid     string `json:"cid"`
	Region  string `json:"region"`
	Origins string `json:"origins"`
	Status  int    `json:"status"`
}

type RepInfo struct {
	Region     string `json:"region"`
	VirtualRep int    `json:"virtual_rep"`
	RealRep    int    `json:"real_rep"`
	MinRep     int    `json:"min_rep"`
	MaxRep     int    `json:"max_rep"`
	Expire     uint64 `json:"expire"`
	Encryption int    `json:"encryption"`
	Status     int    `json:"status"`
}

type TaskResponse struct {
	Fid          string         `json:"fid"` //文件hash值
	Cid          string         `json:"cid"`
	RegionStatus map[string]int `json:"region_status"` //key region, value status
}

type Task struct {
	Fid     string              `bson:"fid,omitempty" json:"fid"`
	Cid     string              `bson:"cid,omitempty" json:"cid"`
	Region  string              `bson:"region,omitempty" json:"region"`   //文件上传区域
	Origins string              `bson:"origins,omitempty" json:"origins"` //文件上传节点
	Reps    map[string]*RepInfo `bson:"reps,omitempty" json:"reps"`       //key region: val: 备份详情
	Status  int                 `bson:"status,omitempty" json:"status"`
}

type ReplicationRequest struct {
	OrderId string           `json:"order_id"`
	Origins string           `json:"origins"`
	Tasks   map[string]*Task `json:"tasks"` //key fid
	Ext     *Extend          `json:"ext,omitempty"`
}

type ReplicationResponse struct {
	Status  int             `json:"status"`
	OrderId string          `json:"order_id"`
	Tasks   []*TaskResponse `json:"tasks"`
}

type TaskReplicateRequest struct {
	OrderId    string            `json:"order_id"`
	Cid        string            `json:"cid"`
	Origins    string            `json:"origins"`
	RealRep    int               `json:"rep"`
	MinRep     int               `json:"min_rep"`
	MaxRep     int               `json:"max_rep"`
	Expire     uint64            `json:"expire"`
	Encryption int               `json:"encryption"`
	NasList    []string          `json:"nas_list,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
	Ext        *Extend           `json:"ext,omitempty"`
}

type TaskReplicateResponse struct {
	Status int `json:"status"`
}

type DeleteOrderRequest struct {
	OrderId string                 `json:"order_id"`
	Tasks   map[string]*UploadTask `json:"tasks"`
	Ext     *Extend                `json:"ext,omitempty"`
}

type DeleteOrderResponse struct {
	Status  int           `json:"status"`
	OrderId string        `json:"order_id"`
	Tasks   []*UploadTask `json:"tasks"`
}

type ChargeRequest struct {
	OrderId   string  `json:"order_id"`
	OrderType int32   `json:"order_type"`
	Tasks     []*Task `json:"tasks"`
	Ext       *Extend `json:"ext,omitempty"`
}

type ChargeResponse struct {
	Status int `json:"status"`
}

type DeleteRequest struct {
	OrderId string  `json:"order_id"`
	Origins string  `json:"origins"`
	Cid     string  `json:"cid"`
	Ext     *Extend `json:"ext,omitempty"`
}

type DeleteResponse struct {
	Status int `json:"status"`
}

type RegionRepResponse struct {
	CidRep map[string]*CidRepInfo `json:"Data"`
	Status int                    `json:"Status"`
}

type CidRepInfo struct {
	PinCount int    `json:"PinCount"` //当前备份数,
	Status   int    `json:"Status"`   //是否备份成功， 0:备份失败,重新备份, 1: 备份成功.
	Message  string `json:"Message"`
}

type RegionRep struct {
	Region      string `json:"region"`
	CurRep      int    `json:"curRep"`      //当前备份数,
	Status      int    `json:"status"`      //是否备份成功， 0:备份失败,重新备份, 1: 备份成功.
	CheckStatus int    `json:"checkStatus"` //校验是否备份成功。 0:校验备份失败， 1: 校验备份成功
}

type TaskRepInfo struct {
	Fid     string                `json:"fid"`
	Cid     string                `json:"cid"`
	Regions map[string]*RegionRep `json:"regions"`
}

type GetOrderRepResponse struct {
	Status  int                     `json:"status"`
	OrderId string                  `json:"order_id"`
	Tasks   map[string]*TaskRepInfo `json:"tasks"`
}

/*type TaskRequest struct {
	Fid     string   `json:"fid"`
	Cid     string   `json:"cid"`
	Regions []string `json:"regions"`
}

type GetOrderRepRequest struct {
	OrderId string         `json:"order_id"`
	Tasks   []*TaskRequest `json:"tasks"`
}*/

type RepTask struct {
	Fid     string   `json:"fid"`
	Cid     string   `json:"cid"`
	Regions []string `json:"regions"`
}

type UploadFinishOrder struct {
	OrderId string     `json:"order_id"`
	Tasks   []*RepTask `json:"tasks"`
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

type Extend struct {
	Ctx context.Context `json:"ctx,omitempty"`
}
