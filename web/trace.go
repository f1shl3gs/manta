package web

import "go.uber.org/zap"

type TraceHandler struct {
	*Router

	logger *zap.Logger
}
