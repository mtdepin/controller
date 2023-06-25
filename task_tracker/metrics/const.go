package metrics

const (
	CHANAL_SIZE = 1024 * 100
)

type Order struct {
	OrderId    string `json:"order_id"`
	OrderType  int32  `json:"order_type"`
	Status     int    `json:"status"`
	UpdateTime int64  `json:"update_time"`
}
