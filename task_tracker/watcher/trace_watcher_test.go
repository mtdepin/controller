package watcher

import (
	"controller/task_tracker/metrics"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTraceWatcher(t *testing.T) {
	GlobalTraceWatcher.Init(1024)
	num := 10
	for i := 1; i < 10; i++ {
		for j := 0; j < num; j++ {
			order := &metrics.Order{
				OrderId:    uuid.NewV4().String(),
				OrderType:  int32(j%2) + 1,
				Status:     i,
				UpdateTime: time.Now().UnixMicro(),
			}

			GlobalTraceWatcher.Add(order)
			GlobalTraceWatcher.Add(order) //模拟重复.
			GlobalTraceWatcher.Delete(order.OrderId)
		}

		time.Sleep(3 * time.Second)
	}

	assert.Equal(t, 0, GlobalTraceWatcher.GetOrderNum())
	time.Sleep(10 * time.Second)
}
