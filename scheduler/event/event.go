package event

const (
	REPLICATE = iota + 1
	DELETE
	CHARGE
	SEARCHREP
)

type Event struct {
	Type int32 //事件类型
	Data interface{}
	Ret  chan interface{}
}
