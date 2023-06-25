package services

import (
	"controller/api"
	"controller/business/database"
	"controller/business/param"
	"controller/business/processor"
)

type Service struct {
	searcher *processor.Searcher
	order    *processor.Order
}

func (p *Service) Init(db *database.DataBase) {
	p.searcher = new(processor.Searcher)
	p.order = new(processor.Order)
	p.order.Init(db)
}

func (p *Service) InitService(order *processor.Order, searcher *processor.Searcher) {
	p.searcher = searcher
	p.order = order
}

func (p *Service) Search(request *param.SearchFileRequest) ([]byte, error) {
	return p.searcher.Search(request)
}

func (p *Service) UploadTask(request *api.UploadTaskRequest) (interface{}, error) {
	return p.order.CreateUploadOrder(request)
}

func (p *Service) UploadFinish(request *param.UploadFinishRequest) (interface{}, error) {
	return p.order.UploadFinish(request)
}

func (p *Service) DownloadTask(request *param.DownloadTaskRequest) (interface{}, error) {
	return p.order.CreateDownloadOrder(request)
}

func (p *Service) DownloadFinish(request *param.DownloadFinishRequest) (interface{}, error) {
	return p.order.DownloadFinish(request)
}

func (p *Service) DeleteFid(request *param.DeleteFidRequest) (interface{}, error) {
	return p.order.DeleteFid(request)
}

/*func (p *Service) GetUploadKNodes(request *kepler_api.GetKNodesRequest) (interface{}, error) {
	return p.order.GetUploadKNodes(request)
}*/

func (p *Service) UploadPieceFid(request *api.UploadPieceFidRequest) (interface{}, error) {
	return p.order.UploadPieceFid(request)
}
