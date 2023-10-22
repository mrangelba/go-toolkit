package events

import (
	"sync"

	dispatcher "github.com/mrangelba/go-toolkit/events/event_dispatcher"
)

var once sync.Once
var instance *dispatcher.EventDispatcher

func GetDispatcher() *dispatcher.EventDispatcher {
	once.Do(func() {
		instance = new()
	})

	return instance
}

func new() *dispatcher.EventDispatcher {
	eventDispatcher := dispatcher.NewEventDispatcher()

	return eventDispatcher
}
