package processor

import (
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"time"
)

//generate order id, 后期采用新算法，生成order_id
func CreateOrderId() string {
	return uuid.NewV4().String()
}

func generateRandNum() int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(10000000)
}
