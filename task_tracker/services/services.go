package services

import (
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

func (p *Service) CreateTask(request *param.CreateTaskRequest) (interface{}, error) {
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
