package param

import "controller/business/dict"

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
	FileName string `json:"file_name"`
	OrderId  string `json:"order_id"`
}

type SearchFileResponse struct {
	Status int              `json:"status"` //1.sucess, 0.fail
	Datas  []*dict.TaskInfo `json:"datas"`
}

type UploadTaskRequest struct {
	RequestId  string             `json:"request_id"`
	UserId     string             `json:"user_id"`
	UploadType int                `json:"upload_type"`
	Tasks      []*dict.UploadTask `json:"tasks"`
	Group      string             `json:"group"`
	NasList    []string           `json:"nas_list,omitempty"`
}

type OrderTaskResponse struct {
	OrderId string `json:"order_id"`
	Status  int    `json:"status"`
}

type NodeListRequst struct {
	Group string `json:"group"`
	Tag   string `json:"tag"`
}

type CreateStrategyRequest struct {
	RequestId string `json:"request_id"`
	OrderId   string `json:"order_id"`
}

type CreateStrategyResponse struct {
	Status int `json:"status"`
}

type UploadTaskResponse struct {
	Status   int      `json:"status"`
	OrderId  string   `json:"order_id"`
	NodeList []string `json:"node_list"`
}

type CheckBalanceRequest struct {
	UserId string `json:"user_id"`
}

type CheckBalanceResponse struct {
	Status int  `json:"status"`
	Enough bool `json:"enough"`
}

type NodeListResponse struct {
	Knodes []string `json:"knodes"`
	Status int      `json:"status"`
}

type CreateTaskRequest struct {
	RequestId string `json:"request_id"`
	Type      int    `json:"task_type"`
}

/*type UploadFinishRequest struct {
	OrderId string            `json:"order_id"`
	Files   map[string]string `json:"files"`
	Status  int               `json:"status"`
}*/

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

type DownloadTaskRequest struct {
	RequestId    string               `json:"request_id"`
	UserId       string               `json:"user_id"`
	Group        string               `json:"group"`
	DownloadType int                  `json:"download_type"`
	Tasks        []*dict.DownloadTask `json:"tasks"`
}

type DownloadTaskResponse struct {
	OrderId string             `json:"order_id"`
	Nodes   map[string][]*Node `json:"nodes"` //下载文件地址
	Status  int                `json:"status"`
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

type Node struct {
	Address string `json:"Address"`
	Weight  int    `json:"Weight"`
}

type CidNodes struct {
	Nodes map[string][]*Node //key: cid , value nodes
}

type GetDownloadNodeResponse struct {
	Data   map[string][]*Node `json:"Data"`
	Status int                `json:"Status"`
}

type UserInfoRequest struct {
	UserId string `json:"user_id"`
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
