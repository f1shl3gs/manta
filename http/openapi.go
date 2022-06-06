package http

import (
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type OpenAPIHandler struct {
	*httprouter.Router

	logger *zap.Logger
}

func NewOpenAPIHandler(backend *Backend, logger *zap.Logger) *OpenAPIHandler {
	panic("todo")
}
