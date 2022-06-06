package log

import (
	"go.uber.org/zap/zapcore"
)

const (
	Console = "console"
	JSON    = "json"
	Logfmt  = "logfmt"
)

type Config struct {
	Format string `json:"format"`
	Level  zapcore.Level
}

func NewConfig() Config {
	return Config{
		Format: "",
		Level:  zapcore.InfoLevel,
	}
}
