package metrics

import (
	"controller/pkg/collector"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var OrderChannel = make(chan *Order, CHANAL_SIZE)

type OrderCollector struct {
	infoDesc *prometheus.Desc
}

func NewOrderCollector() (collector.Collector, error) {
	return &OrderCollector{
		infoDesc: prometheus.NewDesc(
			prometheus.BuildFQName("kepler", "", "order"),
			"kepler upload order ",
			[]string{"order_id", "order_type", "status"},
			nil,
		),
	}, nil
}

func (c *OrderCollector) Update(ch chan<- prometheus.Metric) error {
	nLen := len(OrderChannel)
	if nLen > cap(ch) {
		nLen = cap(ch) - 1
	}

	for i := 0; i < nLen; i++ {
		order := <-OrderChannel
		ch <- prometheus.MustNewConstMetric(
			c.infoDesc,
			prometheus.GaugeValue,
			float64(order.UpdateTime),
			order.OrderId,
			strconv.Itoa(int(order.OrderType)),
			strconv.Itoa(order.Status))
		//logger.Infof("-- collector prometheus: %v, nLen: %v, i: %v", *order, nLen, i)
	}

	return nil
}
