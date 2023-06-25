package services

import "controller/scheduler/event"

func PackageEvent(eventType int32, data interface{}) *event.Event {
	return &event.Event{
		Type: eventType,
		Data: data,
		Ret:  make(chan interface{}),
	}
}
