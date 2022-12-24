package http

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/middlewares"
	"github.com/f1shl3gs/manta/multitsdb"
)

const (
	apiV1Prefix = "/api/v1"
)

type Flusher interface {
	Flush(ctx context.Context) error
}

type Backend struct {
	Flusher Flusher

	router *Router

	BackupService        manta.BackupService
	OrganizationService  manta.OrganizationService
	DashboardService     manta.DashboardService
	UserService          manta.UserService
	PasswordService      manta.PasswordService
	AuthorizationService manta.AuthorizationService
	SessionService       manta.SessionService
	OnBoardingService    manta.OnBoardingService
	CheckService         manta.CheckService
	TaskService          manta.TaskService
	ConfigurationService manta.ConfigurationService
	ScraperTargetService manta.ScraperTargetService
	RegistryService      manta.RegistryService

	TenantStorage         multitsdb.TenantStorage
	TenantTargetRetriever multitsdb.TenantTargetRetriever
}

type Service struct {
	apiHandler    http.Handler
	docHandler    http.Handler
	metricHandler http.Handler
	assetsHandler http.Handler
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if strings.HasPrefix(path, "/api") || path == "/debug/flush" {
		s.apiHandler.ServeHTTP(w, r)
		return
	}

	if path == "/metrics" {
		s.metricHandler.ServeHTTP(w, r)
		return
	}

	if strings.HasPrefix(path, "/docs") {
		s.docHandler.ServeHTTP(w, r)
		return
	}

	if strings.HasPrefix(path, "/debug") {
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}

	// assets handler
	s.assetsHandler.ServeHTTP(w, r)
}

func New(logger *zap.Logger, backend *Backend) *Service {
	// assets handler
	assetsHandler, err := NewAssetsHandler(logger)
	if err != nil {
		panic(err)
	}

	// build api handler
	backend.router = &Router{
		Router: httprouter.New(),
		logger: logger,
	}

	NewOrganizationHandler(backend, logger)
	NewSetupHandler(backend, logger)
	NewSessionHandler(backend.router, logger, backend.UserService, backend.PasswordService, backend.SessionService)
	NewFlushHandler(logger, backend.router, backend.Flusher)
	NewDashboardsHandler(backend, logger)
	NewUserHandler(backend, logger)
	NewConfigurationService(backend, logger)
	NewPromAPIHandler(backend, logger)
	NewScrapeHandler(backend, logger)
	NewRegistryService(backend, logger)
	NewChecksHandler(logger, backend.router, backend.CheckService, backend.TaskService)
	NewTaskHandler(backend, logger)

	ah := &AuthenticationHandler{
		logger:               logger.With(zap.String("handler", "authentication")),
		AuthorizationService: backend.AuthorizationService,
		UserService:          backend.UserService,
		SessionService:       backend.SessionService,
		noAuthRouter:         httprouter.New(),
		handler:              backend.router,
		errorHandler:         backend.router,
	}

	ah.RegisterNoAuthRoute(http.MethodPost, setupPath)
	ah.RegisterNoAuthRoute(http.MethodGet, setupPath)
	ah.RegisterNoAuthRoute(http.MethodPost, signinPath)
	// ah.RegisterNoAuthRoute(http.MethodGet, debugFlushPath)
	ah.RegisterNoAuthRoute(http.MethodGet, "/")
	ah.RegisterNoAuthRoute(http.MethodGet, "/debug/*wild")
	// TODO: add auth in the future
	ah.RegisterNoAuthRoute(http.MethodGet, configurationWithID)
	ah.RegisterNoAuthRoute(http.MethodPost, registryPrefix)

	// set kinds of global middleware
	handler := http.Handler(ah)
	handler = middlewares.Trace(handler)

	// enable access log middleware
	if logger.Core().Enabled(zapcore.DebugLevel) {
		handler = middlewares.Logging(logger, handler)
	}

	return &Service{
		apiHandler:    handler,
		docHandler:    Redoc(),
		metricHandler: promhttp.Handler(),
		assetsHandler: assetsHandler,
	}
}
