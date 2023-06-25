package param

import (
	"context"
	"controller/business/dict"
)

const (
	FAIL        = 0
	SUCCESS     = 1
	INVALID_CID = 2
)

const (
	HaveDevice = 1
	NoDevice   = 2
)

const (
	UPLOAD   = 1
	DOWNLOAD = 2
)

type SearchFileRequest struct {
	FileName string  `json:"file_name"`
	OrderId  string  `json:"order_id"`
	Ext      *Extend `json:"ext,omitempty"`
}

type SearchFileResponse struct {
	Status int              `json:"status"` //1.sucess, 0.fail
	Datas  []*dict.TaskInfo `json:"datas"`
}

type Extend struct {
	Ctx context.Context `json:"ctx,omitempty"`
}

type CreateStrategyRequest struct {
	RequestId string  `json:"request_id"`
	OrderId   string  `json:"order_id"`
	Region    string  `json:"region"`
	Ext       *Extend `json:"ext,omitempty"`
}

type CreateStrategyResponse struct {
	Status int `json:"status"`
}

type UploadTaskResponse struct {
	Status   int          `json:"status"`
	OrderId  string       `json:"order_id"`
	NodeList []*dict.Node `json:"node_list"`
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

type DownloadTaskRequest struct {
	RequestId    string               `json:"request_id"`
	UserId       string               `json:"user_id"`
	Group        string               `json:"group"`
	DownloadType int                  `json:"download_type"`
	Tasks        []*dict.DownloadTask `json:"tasks"`
	Ext          *Extend              `json:"ext,omitempty"`
}

type DownloadTaskResponse struct {
	OrderId string                  `json:"order_id"`
	Nodes   map[string][]*dict.Node `json:"nodes"` //下载文件地址
	Status  int                     `json:"status"`
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

type CidNodes struct {
	Nodes map[string][]*dict.Node //key: cid , value nodes
}

type GetDownloadNodeResponse struct {
	Data   map[string][]*dict.Node `json:"Data"`
	Status int                     `json:"Status"`
}

type UserInfoRequest struct {
	UserId string  `json:"user_id"`
	Ext    *Extend `json:"ext,omitempty"`
}

type UserInfoResponse struct {
	UserId   string `json:"user_id"`
	UserNick string `json:"user_nick"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	UnitName string `json:"unit_name"`
	Secret   string `json:"secret"`
	Pris     []int  `json:"pris"`
	Status   int    `json:"status"`
}

//orderId, fid.
type DeleteFidRequest struct {
	RequestId string          `json:"request_id"`
	UserId    string          `json:"user_id"`
	OrderId   string          `json:"order_id"`
	Fids      map[string]bool `json:"fids"`
	Ext       *Extend         `json:"ext,omitempty"`
}

type DeleteFidResponse struct {
	OrderId string         `json:"order_id"`
	Fids    map[string]int `json:"fids"` //每个fid的删除状态
	Status  int            `json:"status"`
}
