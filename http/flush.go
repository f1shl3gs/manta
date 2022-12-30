package http

import (
	"net/http"

	"go.uber.org/zap"
)

const debugFlushPath = "/debug/flush"

func NewFlushHandler(logger *zap.Logger, backend *Backend) {
	flusher := backend.Flusher
	if flusher == nil {
		return
	}

	router := backend.router
	router.HandlerFunc(http.MethodGet, debugFlushPath, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := flusher.Flush(ctx)
		if err != nil {
			router.HandleHTTPError(ctx, err, w)
			return
		}

		logger.Info("Flush success")
	})
}
