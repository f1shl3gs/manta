package http

import (
	"net/http"

	"go.uber.org/zap"
)

const DebugFlushPath = "/debug/flush"

func NewFlushHandler(logger *zap.Logger, router *Router, flusher Flusher) {
	if flusher == nil {
		return
	}

	router.HandlerFunc(http.MethodGet, DebugFlushPath, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := flusher.Flush(ctx)
		if err != nil {
			router.HandleHTTPError(ctx, err, w)
			return
		}

		logger.Info("Flush success")
	})
}
