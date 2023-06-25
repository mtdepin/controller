package event

import "controller/strategy/dict"

type Event struct {
	//Type  int //eventType
	OrderId string
	Data    *dict.Task
	Ret     chan int
}
