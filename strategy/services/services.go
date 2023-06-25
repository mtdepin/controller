package services

import (
	"controller/strategy/database"
	"controller/strategy/param"
	"controller/strategy/processor"
)

type Service struct {
	strategy *processor.Strategy
}

func (p *Service) Init(db *database.DataBase) {
	p.strategy = new(processor.Strategy)
	p.strategy.Init(db, 300)
}

func (p *Service) CreateStrategy(request *param.CreateStrategyRequest) (interface{}, error) {
	return p.strategy.CreateStrategy(request)
}

func (p *Service) GetReplicateStrategy(request *param.GetStrategyRequset) (interface{}, error) {
	return p.strategy.GetReplicateStrategy(request)
}

func (p *Service) GetOrderDeleteStrategy(request *param.GetStrategyRequset) (interface{}, error) {
	return p.strategy.GetOrderDeleteStrategy(request)
}

func (p *Service) GetFidDeleteStrategy(request *param.GetFidDeleteStrategyRequest) (interface{}, error) {
	return p.strategy.GetFidDeleteStrategy(request)
}
