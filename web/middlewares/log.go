package middlewares

import (
	"net/http"
	"time"

	"github.com/f1shl3gs/manta/pkg/tracing"
	"go.uber.org/zap"
)

type logResponse struct {
	http.ResponseWriter

	status int
}

func (l *logResponse) WriteHeader(statusCode int) {
	l.status = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func Log(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &logResponse{
			ResponseWriter: w,
		}

		start := time.Now()
		next.ServeHTTP(lw, r)
		latency := time.Since(start)

		id, _, _ := tracing.InfoFromContext(r.Context())

		logger.Info("serve http request",
			zap.String("trace_id", id),
			zap.String("remote", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Int("status", lw.status),
			zap.Duration("latency", latency))
	})
}
