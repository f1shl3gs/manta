package middlewares

import (
	"net/http"
	"time"

	"github.com/f1shl3gs/manta/pkg/tracing"
	"go.uber.org/zap"
)

func Log(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newRecordableResponse(w)

		start := time.Now()
		next.ServeHTTP(rw, r)
		latency := time.Since(start)

		id, _, _ := tracing.InfoFromContext(r.Context())

		logger.Info("serve http request",
			zap.String("trace_id", id),
			zap.String("remote", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Int("status", rw.Status()),
			zap.Int("written", rw.Written()),
			zap.Duration("latency", latency))
	})
}
