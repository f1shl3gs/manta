package web

import (
	"net/http"
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
	ua "github.com/mileusna/useragent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/web/middlewares"
)

type Backend struct {
	manta.HTTPErrorHandler

	OtclService          manta.OtclService
	BackupService        manta.BackupService
	NodeService          manta.NodeService
	OrganizationService  manta.OrganizationService
	CheckService         manta.CheckService
	TaskService          manta.TaskService
	DatasourceService    manta.DatasourceService
	TemplateService      manta.TemplateService
	UserService          manta.UserService
	AuthorizationService manta.AuthorizationService
}

func New(logger *zap.Logger, backend *Backend) http.Handler {
	router := NewRouter()

	// healthz
	router.Handler(http.MethodGet, "/healthz", newHealthzHandler(logger))

	otclService(logger, router, backend)

	// readiness
	router.Handler(http.MethodGet, "/ready", ReadyHandler())

	// organizations
	NewOrganizationHandler(logger, router, backend)

	{
		// prometheus
		mh := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			MaxRequestsInFlight: 3,
		})

		router.Handler(http.MethodGet, "/metrics", mh)
	}

	{
		mux := &http.ServeMux{}
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/traces", pprof.Trace)

		// pprof
		router.GET("/debug/pprof/*dummy", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			mux.ServeHTTP(w, r)
		})
	}

	// tracing
	h := middlewares.Log(logger, router)
	h = Trace(h)

	return h
}

func Trace(next http.Handler) http.Handler {
	name := "manta"
	fn := func(w http.ResponseWriter, r *http.Request) {
		span, r := tracing.ExtractFromHTTPRequest(r, name)
		defer span.Finish()

		span.LogKV("user_agent", UserAgent(r))
		for k, v := range r.Header {
			if len(v) == 0 {
				continue
			}

			if k == "Authorization" || k == "User-Agent" {
				continue
			}

			// If header has multiple values, only the first value will be logged on the traces.
			span.LogKV(k, v[0])
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func UserAgent(r *http.Request) string {
	header := r.Header.Get("User-Agent")
	if header == "" {
		return "unknown"
	}

	return ua.Parse(header).Name
}
