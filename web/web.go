package web

import (
	"github.com/f1shl3gs/manta/authorization"
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
	OrganizationService  manta.OrganizationService
	CheckService         manta.CheckService
	TaskService          manta.TaskService
	DatasourceService    manta.DatasourceService
	DashboardService     manta.DashboardService
	TemplateService      manta.TemplateService
	UserService          manta.UserService
	PasswordService      manta.PasswordService
	AuthorizationService manta.AuthorizationService
	Keyring              manta.Keyring
	SessionService       manta.SessionService
}

func New(logger *zap.Logger, backend *Backend) http.Handler {
	router := NewRouter()

	// static
	fileServer := http.FileServer(http.FS(manta.Assets))
	router.GET("/ui/*filepath", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		r.URL.Path = manta.UIPrefix + params.ByName("filepath")
		fileServer.ServeHTTP(w, r)
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

	userService(logger, router, backend.UserService)

	NewSetupHandler(router, logger, backend)

	NewSessionHandler(router, logger, backend.UserService, backend.PasswordService, backend.SessionService)

	NewKeyringHandler(router, logger, backend.Keyring)

	// datasource
	DatasourceService(logger, router, backend.DatasourceService)

	// dashboard
	dh := &DashboardHandler{
		Router:           router,
		logger:           logger,
		dashboardService: backend.DashboardService,
	}

	NewDashboardService(dh)

	// and more

	// tracing
	h := middlewares.Log(logger, router)
	h = Trace(h)

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
