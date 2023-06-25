package watcher

import (
	"controller/pkg/collector"
	"controller/task_tracker/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

var GlobalTraceWatcher = new(TraceWatcher)

type TraceWatcher struct {
	orders map[string]*metrics.Order
	mutex  *sync.RWMutex
}

func (p *TraceWatcher) Init(nSize int) {
	p.orders = make(map[string]*metrics.Order, nSize)
	p.mutex = new(sync.RWMutex)
	p.InitRegister()
}

func (p *TraceWatcher) InitRegister() {
	collector.RegisterCollector("OrderCollector", metrics.NewOrderCollector)

	nodeCollector, err := collector.NewNodeCollector()
	if err != nil {
		return
	}

	r := prometheus.NewRegistry()
	if err = r.Register(nodeCollector); err != nil {
		return
	}

	registry := prometheus.NewRegistry()
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{registry, r, prometheus.DefaultGatherer},
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
			Registry:      registry,
		},
	)

	http.Handle("/metrics", handler)
	go func() {
		http.ListenAndServe("0.0.0.0:8989", nil)
	}()
}

func (p *TraceWatcher) Add(order *metrics.Order) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if value, ok := p.orders[order.OrderId]; ok {
		if value.Status == order.Status { //过滤
			return
		} else { //不等, 新值
			value.Status = order.Status
			value.UpdateTime = order.UpdateTime
		}
	} else { //不存在，添加.
		newOrder := *order
		p.orders[order.OrderId] = &newOrder
	}

	metrics.OrderChannel <- order
}

func (p *TraceWatcher) Delete(orderId string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.orders, orderId)
}

func (p *TraceWatcher) GetOrderNum() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return len(p.orders)
}
