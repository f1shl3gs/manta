package log

import (
	"context"

	"go.uber.org/zap"
)

type key string

const (
	loggerKey key = "logger"
)

// NewContextWithLogger returns a new context with file added.
func NewContextWithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// LoggerFromContext returns the zap.Logger associated with ctx or nil if no logger has been assigned.
func LoggerFromContext(ctx context.Context) *zap.Logger {
	l, _ := ctx.Value(loggerKey).(*zap.Logger)
	return l
}
