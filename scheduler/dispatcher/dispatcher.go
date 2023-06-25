package dispatcher

import (
	e "controller/scheduler/event"
	"controller/scheduler/manager"
)

type Dispatcher struct {
	handers []func(event *e.Event)
}

func (p *Dispatcher) Init(manager *manager.Manager) {
	//p.pipeline = make(chan *e.Event, size)
	p.handers = make([]func(event *e.Event), 10, 10)

	p.handers[e.REPLICATE] = manager.GetHandler(e.REPLICATE)
	p.handers[e.CHARGE] = manager.GetHandler(e.CHARGE)
	p.handers[e.DELETE] = manager.GetHandler(e.DELETE)
	p.handers[e.SEARCHREP] = manager.GetHandler(e.SEARCHREP)
	p.handers[e.DELETEFID] = manager.GetHandler(e.DELETEFID)

	//go p.Dispatch()
}

func (p *Dispatcher) AddReplicateEvent(event *e.Event) {
	p.handers[e.REPLICATE](event)
}

func (p *Dispatcher) AddDeleteEvent(event *e.Event) {
	p.handers[e.DELETE](event)
}

func (p *Dispatcher) AddChargeEvent(event *e.Event) {
	p.handers[e.CHARGE](event)
}

func (p *Dispatcher) AddSearchRepEvent(event *e.Event) {
	p.handers[e.SEARCHREP](event)
}

func (p *Dispatcher) AddDeleteFidEvent(event *e.Event) {
	p.handers[e.DELETEFID](event)
}
