package dispatcher

import (
	"context"
	"sync"
)

type EventHandler interface {
	Handle(ctx context.Context, payload interface{}, wg *sync.WaitGroup)
}

type Dispatcher interface {
	Register(eventName string, handler EventHandler) error
	Dispatch(ctx context.Context, event string, payload interface{})
	Remove(eventName string, handler EventHandler)
	Has(eventName string, handler EventHandler) bool
	Clear()
}
