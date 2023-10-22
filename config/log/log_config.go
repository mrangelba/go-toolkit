package log

import (
	"os"
)

type LogConfig struct {
	Level string
}

func Get() LogConfig {
	return LogConfig{
		Level: os.Getenv("LOG_LEVEL"),
	}
}
