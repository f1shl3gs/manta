package router

import (
    "net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
)

func Metrics() Middleware {
    var (
        namespace = "manta"
        subsystem = "http"
        labels    = []string{"method", "handler", "code"}
    )

    requestsLatency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Namespace: namespace,
        Subsystem: subsystem,
        Name:      "requests_latency_seconds",
        Help:      "Histogram of times spent for end-to-end latency",
        Buckets:   prometheus.ExponentialBuckets(1e-3, 5, 7),
    }, labels)

    prometheus.DefaultRegisterer.MustRegister(requestsLatency)

    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            rw := newRecordableResponse(w)

            start := time.Now()
            next.ServeHTTP(rw, r)
            latency := time.Since(start)

            params := httprouter.ParamsFromContext(r.Context())
            labelValues := []string{r.Method, params.MatchedRoutePath(), strconv.Itoa(rw.Status())}

            requestsLatency.WithLabelValues(labelValues...).Observe(latency.Seconds())
        }
    }
}
