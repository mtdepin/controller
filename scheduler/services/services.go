package services

import (
	"controller/scheduler/database"
	"controller/scheduler/dispatcher"
	"controller/scheduler/event"
	"controller/scheduler/manager"
	"controller/scheduler/param"
)

const (
	MSG_QUEEN_SIZE = 1024
	PROCESSOR_NUM  = 20
)

type Service struct {
	manager    *manager.Manager
	dispatcher *dispatcher.Dispatcher
}

func (p *Service) Init(db *database.DataBase) {
	p.manager = new(manager.Manager)
	p.manager.Init(MSG_QUEEN_SIZE, PROCESSOR_NUM, db)

	p.dispatcher = new(dispatcher.Dispatcher)
	p.dispatcher.Init(p.manager)
}

func (p *Service) Replicate(request *param.ReplicationRequest) (interface{}, error) {
	event := PackageEvent(event.REPLICATE, request)
	p.dispatcher.AddReplicateEvent(event)
	return <-event.Ret, nil
}

func (p *Service) DeleteOrder(request *param.DeleteOrderRequest) (interface{}, error) {
	event := PackageEvent(event.DELETE, request)
	p.dispatcher.AddDeleteEvent(event)
	return <-event.Ret, nil
}

func (p *Service) Charge(request *param.ChargeRequest) (interface{}, error) {
	event := PackageEvent(event.CHARGE, request)
	p.dispatcher.AddChargeEvent(event)
	return <-event.Ret, nil
}

func (p *Service) SearchRep(request *param.UploadFinishOrder) (interface{}, error) {
	event := PackageEvent(event.SEARCHREP, request)
	p.dispatcher.AddSearchRepEvent(event)
	return <-event.Ret, nil
}

func (p *Service) DeleteFid(request *param.DeleteOrderFidRequest) (interface{}, error) {
	event := PackageEvent(event.DELETE, request)
	p.dispatcher.AddDeleteFidEvent(event)
	return <-event.Ret, nil
}
