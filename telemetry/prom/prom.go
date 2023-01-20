package prom

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Registry struct {
	*prometheus.Registry

	logger *zap.Logger
}

func NewRegistry(logger *zap.Logger) *Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	return &Registry{
		Registry: reg,
		logger:   logger,
	}
}

// HTTPHandler returns an http.Handler for the registry,
// so that the /metrics HTTP handler is uniformly configured across all apps in the platform.
func (r *Registry) HTTPHandler() http.Handler {
	opts := promhttp.HandlerOpts{
		ErrorLog: promLogger{reg: r},
		// TODO(mr): decide if we want to set MaxRequestsInFlight or Timeout.
	}
	return promhttp.HandlerFor(r.Registry, opts)
}

// promLogger satisfies the promhttp.logger interface with the registry.
// Because normal usage is that WithLogger is called after HTTPHandler,
// we refer to the Registry rather than its logger.
type promLogger struct {
	reg *Registry
}

var _ promhttp.Logger = (*promLogger)(nil)

// Println implements promhttp.logger.
func (pl promLogger) Println(v ...interface{}) {
	pl.reg.logger.Sugar().Info(v...)
}
