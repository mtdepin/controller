package manager

import (
	"controller/scheduler/database"
	"controller/scheduler/dict"
	"controller/scheduler/event"
	"controller/scheduler/processor"
	"fmt"
)

type Manager struct {
	domain             *database.Domain
	chargeProcessor    *processor.ChargeProcessor
	deleteProcessor    *processor.DeleteProcessor
	replicateProcessor *processor.ReplicateProcessor
	searchRepProcessor *processor.SearchRepProcessor
}

func (p *Manager) Init(size, num int32, db *database.DataBase) {
	p.domain = new(database.Domain)
	p.domain.Init(db)

	domainMap, err := p.getDomains()
	if err != nil {
		panic(fmt.Sprintf("Manager Init getDomains fail %v", err.Error()))
	}

	p.chargeProcessor = new(processor.ChargeProcessor)
	p.chargeProcessor.Init(size, num)

	p.deleteProcessor = new(processor.DeleteProcessor)
	p.deleteProcessor.Init(size, num, domainMap)

	p.replicateProcessor = new(processor.ReplicateProcessor)
	p.replicateProcessor.Init(size, num, domainMap)

	p.searchRepProcessor = new(processor.SearchRepProcessor)
	p.searchRepProcessor.Init(size, num, domainMap)

}

func (p *Manager) GetHandler(eventType int32) func(event *event.Event) {
	switch eventType {
	case event.REPLICATE:
		return p.replicateProcessor.AddEvent
	case event.DELETE:
		return p.deleteProcessor.AddEvent
	case event.CHARGE:
		return p.chargeProcessor.AddEvent
	case event.SEARCHREP:
		return p.searchRepProcessor.AddEvent
	default:
		return nil
	}

	return nil
}

func (p *Manager) getDomains() (map[string]*dict.DomainInfo, error) {
	domains, err := p.domain.Load()
	if err != nil {
		return nil, err
	}

	domainMap := make(map[string]*dict.DomainInfo, len(*domains))
	for i, _ := range *domains {
		domainMap[(*domains)[i].Region] = &(*domains)[i]
	}

	return domainMap, nil
}
