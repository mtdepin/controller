package services

import (
	"controller/api"
	"controller/task_tracker/database"
	"controller/task_tracker/manager"
	"controller/task_tracker/param"
)

type Service struct {
	manager *manager.Manager
}

func (p *Service) Init(db *database.DataBase) {
	p.manager = new(manager.Manager)
	p.manager.Init(db)
}

func (p *Service) CreateTask(request *api.CreateTaskRequest) (interface{}, error) {
	return p.manager.CreateTask(request)
}

func (p *Service) UploadFinish(request *param.UploadFinishRequest) (interface{}, error) {
	return p.manager.UploadFinish(request)
}

func (p *Service) CallbackUpload(request *param.CallbackUploadRequest) (interface{}, error) {
	return p.manager.CallbackUpload(request)
}

func (p *Service) CallbackReplicate(request *param.CallbackRepRequest) (interface{}, error) {
	return p.manager.CallbackRep(request)
}

func (p *Service) CallbackDelete(request *param.CallbackDeleteRequest) (interface{}, error) {
	return p.manager.CallbackDelete(request)
}

func (p *Service) CallbackCharge(request *param.CallbackChargeRequest) (interface{}, error) {
	return p.manager.CallbackCharge(request)
}

func (p *Service) DownloadFinish(request *param.DownloadFinishRequest) (interface{}, error) {
	return p.manager.DownloadFinish(request)
}

//重复文件，order_id
func (p *Service) DeleteFid(request *param.DeleteFidRequest) (interface{}, error) {
	return p.manager.DeleteFid(request)
}

func (p *Service) CallbackDownload(request *param.CallbackDownloadRequest) (interface{}, error) {
	return p.manager.CallbackDownload(request)
}

func (p *Service) GetOrderDetail(request *param.SearchOrderRequest) (interface{}, error) {
	return p.manager.GetOrderDetail(request)
}

func (p *Service) GetPageOrders(request *param.OrderPageQueryRequest) (interface{}, error) {
	return p.manager.GetPageOrders(request)
}

func (p *Service) GetFidDetail(request *param.SearchFidRequest) (interface{}, error) {
	return p.manager.GetFidDetail(request)
}

func (p *Service) GetPageFids(request *param.FidPageQueryRequest) (interface{}, error) {
	return p.manager.GetPageFids(request)
}

func (p *Service) UploadPieceFid(request *api.UploadPieceFidRequest) (interface{}, error) {
	return p.manager.UploadPieceFid(request)
}
