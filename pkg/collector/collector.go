package collector

import (
	"controller/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	factories = make(map[string]func() (Collector, error))
)

type Collector interface {
	Update(ch chan<- prometheus.Metric) error
}

type NodeCollector struct {
	Collectors map[string]Collector
}

func RegisterCollector(collector string, factory func() (Collector, error)) {
	factories[collector] = factory
}

func NewNodeCollector() (*NodeCollector, error) {
	collectors := make(map[string]Collector)

	for key, factorie := range factories {
		collector, err := factorie()
		if err != nil {
			return nil, err
		}

		collectors[key] = collector
	}

	return &NodeCollector{Collectors: collectors}, nil
}

func (n NodeCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (n NodeCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.Collectors))
	for name, c := range n.Collectors {
		go func(name string, c Collector) {
			execute(c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(c Collector, ch chan<- prometheus.Metric) {
	if err := c.Update(ch); err != nil {
		logger.Warnf("update metrics fail: %v", err.Error())
	}
}
