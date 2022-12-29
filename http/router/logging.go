package router

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta/pkg/tracing"
)

func Logging(logger *zap.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rw := newRecordableResponse(w)

			start := time.Now()
			next.ServeHTTP(rw, r)
			latency := time.Since(start)

			id, _, _ := tracing.InfoFromContext(r.Context())

			logger.Debug("Serve http request",
				zap.String("trace_id", id),
				zap.String("remote", r.RemoteAddr),
				zap.String("method", r.Method),
				zap.String("url", r.URL.Path),
				zap.Int("status", rw.Status()),
				zap.Int("written", rw.Written()),
				zap.Duration("latency", latency))
		}
	}
}
