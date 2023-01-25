package http

import (
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"

	"go.uber.org/zap"
)

const (
	operationLogPath = apiV1Prefix + `/oplogs/:id`
)

type OperationLogHandler struct {
	*router.Router
	logger *zap.Logger

	oplogService manta.OperationLogService
}

func NewOperationLogHandler(backend *Backend) {
	h := &OperationLogHandler{
		Router:       backend.router,
		oplogService: backend.OperationLogService,
	}

	h.HandlerFunc(http.MethodGet, operationLogPath, h.handleList)
}

func (h *OperationLogHandler) handleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromPath(r)
	if err != nil {
		return
	}

	opts, err := manta.DecodeFindOptions(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	logs, _, err := h.oplogService.FindOperationLogsByUser(ctx, id, opts)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, logs); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func findOplogByResourceID(r *http.Request, oplogService manta.OperationLogService) ([]*manta.OperationLogEntry, int, error) {
	ctx := r.Context()

	id, err := idFromPath(r)
	if err != nil {
		return nil, 0, err
	}

	opts, err := manta.DecodeFindOptions(r)
	if err != nil {
		return nil, 0, err
	}

	return oplogService.FindOperationLogsByID(ctx, id, opts)
}
