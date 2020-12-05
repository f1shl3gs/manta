package web

import "go.uber.org/zap"

const (
	DatasourcePath = "/api/v1/datasource/:id/*path"
)

// todo: datasource cache

type DatasourceHandler struct {
	logger *zap.Logger
}
