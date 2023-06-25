package api

type UploadPieceFidResponse struct {
	OrderId string            `json:"order_id"`
	RepFids map[string]string `json:"rep_fids"` //每个fid的删除状态
	Status  int               `json:"status"`
}

type GetKNodesResponse struct {
	Knodes []*Node `json:"knodes"`
	Status int     `json:"status"`
}

type Node struct {
	Address string `json:"Address"`
	Weight  int    `json:"Weight"`
	RtcId   string `json:"RtcId"`
}

type UploadTaskResponse struct {
	Status   int     `json:"status"`
	OrderId  string  `json:"order_id"`
	NodeList []*Node `json:"node_list"`
	Group    string  `json:"group"`
}

type CreateTaskResponse struct {
	Status  int    `json:"status"`
	OrderId string `json:"order_id"`
}

type OrderTaskResponse struct {
	OrderId string            `json:"order_id"`
	Status  int               `json:"status"`
	Fids    map[string]string `json:"fids"` //key:fid ,value:cid
}

type NodeListResponse struct {
	Knodes []*Node `json:"knodes"`
	Status int     `json:"status"`
}

type CheckBalanceResponse struct {
	Status int  `json:"status"`
	Enough bool `json:"enough"`
}

type GetRmSpaceResponse struct {
	Status int         `json:"status"`
	Region *RegionInfo `json:"region"`
}

type RegionInfo struct {
	RegionName   string `json:"regionName"`
	TotalStorage uint64 `json:"totalStorage"`
	ValidStorage uint64 `json:"validStorage"`
}
