package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func Metrics(reg prometheus.Registerer, next http.Handler) http.Handler {
	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "manta",
		Subsystem: "http",
		Name:      "request_duration_seconds",
	}, []string{"method", "path"})

	status := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "manta",
		Subsystem: "http",
		Name:      "request_total",
	}, []string{"method", "path", "code"})

	inflights := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "manta",
		Subsystem: "http",
		Name:      "inflight_requests",
	})

	reg.MustRegister(latency, status, inflights)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newRecordableResponse(w)

		start := time.Now()
		inflights.Inc()

		next.ServeHTTP(rw, r)

		latency.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
		status.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.Status())).Inc()
		inflights.Dec()
	})
}
