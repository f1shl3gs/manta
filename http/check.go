package http

import (
	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

const (
	ChecksPrefix = `/api/v1/checks`
	ChecksIDPath = `/api/v1/checks/:id`
)

type ChecksHandler struct {
	*Router

	logger       *zap.Logger
	checkService manta.CheckService
	taskService  manta.TaskService
}
