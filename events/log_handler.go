package events

import (
	"context"
	"sync"

	"github.com/mrangelba/go-toolkit/logger"
)

type LogEventHandler struct {
	EventName string
}

func NewLogEventHandler(name string) *LogEventHandler {
	return &LogEventHandler{
		EventName: name,
	}
}

func (e *LogEventHandler) Handle(ctx context.Context, payload interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	log := logger.Get()

	log.Info().Any("payload", payload).Str("event", e.EventName).Msg("Event dispatched")
}
