package http

import (
    "github.com/f1shl3gs/manta/http/router"
    "net/http"

	"github.com/f1shl3gs/manta"

	"go.uber.org/zap"
)

const (
	TaskIDPath     = apiV1Prefix + `/tasks/:id`
	TaskRunsPrefix = TaskIDPath + `/runs`

	TaskDefaultPageSize = 100
	TaskMaxPageSize     = 500
)

type TaskHandler struct {
	*router.Router
	logger *zap.Logger

	taskService manta.TaskService
}

func NewTaskHandler(backend *Backend, logger *zap.Logger) {
	h := &TaskHandler{
		Router:      backend.router,
		logger:      logger.With(zap.String("handler", "task")),
		taskService: backend.TaskService,
	}

	h.HandlerFunc(http.MethodGet, TaskIDPath, h.findTask)
	h.HandlerFunc(http.MethodGet, TaskRunsPrefix, h.findTaskRuns)
}

func (h *TaskHandler) findTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	task, err := h.taskService.FindTaskByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, task); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *TaskHandler) findTaskRuns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	limit, err := limitFromQuery(r, TaskDefaultPageSize, TaskMaxPageSize)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	runs, _, err := h.taskService.FindRuns(ctx, manta.RunFilter{
		Task:   id,
		Limit:  limit,
		After:  nil,
		Before: nil,
	})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, &runs); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
