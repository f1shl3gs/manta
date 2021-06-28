package web

import (
	"context"
	"net/http"
	"net/http/pprof"
	"strings"

	tsdb2 "github.com/f1shl3gs/manta/pkg/tsdb"
	"github.com/julienschmidt/httprouter"
	ua "github.com/mileusna/useragent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorization"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/f1shl3gs/manta/web/middlewares"
)

type Flusher interface {
	Flush(ctx context.Context) error
}

type Backend struct {
	manta.HTTPErrorHandler

	Flusher Flusher

	TenantStorage               tsdb2.TenantStorage
	OtclService                 manta.OtclService
	BackupService               manta.BackupService
	OrganizationService         manta.OrganizationService
	CheckService                manta.CheckService
	TaskService                 manta.TaskService
	DashboardService            manta.DashboardService
	TemplateService             manta.TemplateService
	UserService                 manta.UserService
	PasswordService             manta.PasswordService
	AuthorizationService        manta.AuthorizationService
	Keyring                     manta.Keyring
	SessionService              manta.SessionService
	SecretService               manta.SecretService
	ScrapeService               manta.ScraperTargetService
	VariableService             manta.VariableService
	NotificationEndpointService manta.NotificationEndpointService
}

func New(logger *zap.Logger, backend *Backend, accessLog bool) http.Handler {
	router := NewRouter()

	assetsHandler, err := NewAssetsHandler(logger)
	if err != nil {
		panic(err)
	}

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/api/") {
			http.NotFound(w, r)
			return
		}

		assetsHandler.ServeHTTP(w, r)
	})

	// healthz
	router.Handler(http.MethodGet, "/healthz", newHealthzHandler(logger))

	// readiness
	router.Handler(http.MethodGet, "/ready", ReadyHandler())

	// organizations
	NewOrganizationHandler(logger, router, backend)

	otclService(logger, router, backend)

	{
		// prometheus
		router.HandlerFunc(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			mfs, err := prometheus.DefaultGatherer.Gather()
			if err != nil {
				logger.Error("gathering metrics failed",
					zap.Error(err))
				router.HandleHTTPError(ctx, err, w)
				return
			}

			enc := expfmt.NewEncoder(w, expfmt.FmtText)
			for _, mf := range mfs {
				err = enc.Encode(mf)
				if err != nil {
					logger.Warn("encode metric family failed",
						zap.Stringp("name", mf.Name),
						zap.Error(err))
				}
			}

			if closer, ok := enc.(expfmt.Closer); ok {
				closer.Close()
			}
		})
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

	userService(logger, router, backend.UserService)

	NewSetupHandler(router, logger, backend)

	NewSessionHandler(router, logger, backend.UserService, backend.PasswordService, backend.SessionService)

	NewKeyringHandler(router, logger, backend.Keyring)

	// dashboard
	dh := &DashboardHandler{
		Router:           router,
		logger:           logger,
		dashboardService: backend.DashboardService,
	}

	NewDashboardService(dh)

	NewScrapeHandler(logger, router, backend.ScrapeService)

	NewQueryHandler(logger, router, backend.TenantStorage)

	NewChecksHandler(logger, router, backend.CheckService, backend.TaskService)

	NewNotificationEndpointHandler(logger, router, backend.NotificationEndpointService)

	NewVariableHandler(logger, router, backend.VariableService)

	NewSecretHandler(logger, router, backend.SecretService)

	NewProfileHandler(logger, router)

	// and more

	if backend.Flusher != nil {
		flusher := backend.Flusher
		router.HandlerFunc(http.MethodGet, "/kv/flush", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			err := flusher.Flush(r.Context())
			if err != nil {
				router.HandleHTTPError(ctx, err, w)
				return
			}

			w.WriteHeader(http.StatusOK)
		})
	}

	var h http.Handler = router

	// middlewares
	h = Trace(h)
	h = middlewares.Metrics(prometheus.DefaultRegisterer, h)
	// h = middlewares.Gzip(h)

	// access log
	if accessLog {
		h = middlewares.Log(logger, h)
	} else {
		logger.Debug("Access log is disabled")
	}

	ah := &AuthenticationHandler{
		logger:               logger,
		AuthorizationService: backend.AuthorizationService,
		UserService:          backend.UserService,
		Keyring:              backend.Keyring,
		handler:              h,
		errorHandler:         router,
		noAuthRouter:         httprouter.New(),
		tokenParser:          authorization.NewTokenParser(backend.Keyring),
		SessionService:       backend.SessionService,
	}

	ah.RegisterNoAuthRoute(http.MethodPost, "/api/v1/signin")
	ah.RegisterNoAuthRoute(http.MethodPost, "/api/v1/signout")
	ah.RegisterNoAuthRoute(http.MethodPost, "/api/v1/setup")
	ah.RegisterNoAuthRoute(http.MethodGet, "/metrics")
	ah.RegisterNoAuthRoute(http.MethodGet, "/debug/pprof/*all")
	ah.RegisterNoAuthRoute(http.MethodGet, "/debug/pprof")
	ah.RegisterNoAuthRoute(http.MethodGet, "/kv/flush")
	ah.RegisterNoAuthRoute(http.MethodGet, "/")

	return ah
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
