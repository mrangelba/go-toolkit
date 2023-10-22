package logger

import (
	"strings"
	"sync"

	"github.com/mrangelba/go-toolkit/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var once sync.Once
var instance Logger

func Get() Logger {
	once.Do(func() {
		instance = new()
	})

	return instance
}

func new() Logger {
	cfg := config.Get()

	l := log.With().Str("service", strings.ToLower(cfg.Service.Name))

	Configure(
		Options{
			LogLevel: cfg.Log.Level,
			JSON:     cfg.IsProd(),
		},
	)

	if !DefaultOptions.Concise && len(DefaultOptions.Tags) > 0 {
		l = l.Fields(map[string]interface{}{
			"tags": DefaultOptions.Tags,
		})
	}

	return Logger{l.Logger()}
}

type Logger struct {
	zerolog.Logger
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatal().Msgf(format, v...)
}

func (l Logger) Errorf(format string, v ...interface{}) {
	l.Logger.Error().Msgf(format, v...)
}

func (l Logger) Warnf(format string, v ...interface{}) {
	l.Logger.Warn().Msgf(format, v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	l.Logger.Info().Msgf(format, v...)
}

func (l Logger) Debugf(format string, v ...interface{}) {
	l.Logger.Debug().Msgf(format, v...)
}

func (l Logger) Tracef(format string, v ...interface{}) {
	l.Logger.Trace().Msgf(format, v...)
}
