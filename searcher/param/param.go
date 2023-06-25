package param

import "controller/searcher/dict"

const (
	SUCCESS = 1
	FAIL    = 0
)

type SearchFileRequest struct {
	FileName string `json:"file_name"`
	OrderId  string `json:"order_id"`
}

type SearchFileResponse struct {
	Status int              `json:"status"` //1.sucess, 0.fail
	Datas  []*dict.TaskInfo `json:"datas"`
}
