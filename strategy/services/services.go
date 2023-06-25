package services

import (
	"controller/strategy/database"
	"controller/strategy/param"
	"controller/strategy/processor"
)

const (
	FIDREP_QUEEN_SIZE = 1024
	PROCESSOR_NUM     = 1
)

type Service struct {
	strategy *processor.Strategy
}

func (p *Service) Init(db *database.DataBase) {
	p.strategy = new(processor.Strategy)
	p.strategy.Init(db, FIDREP_QUEEN_SIZE, PROCESSOR_NUM)
}

func (p *Service) CreateStrategy(request *param.CreateStrategyRequest) (interface{}, error) {
	return p.strategy.CreateStrategy(request)
}

func (p *Service) GetReplicateStrategy(request *param.GetStrategyRequset) (interface{}, error) {
	return p.strategy.GetStrategy(request)
}

func (p *Service) GetDeleteReplicateStrategy(request *param.GetDeleteStrategyRequset) (interface{}, error) {
	return p.strategy.GetDeleteStrategy(request)
}
