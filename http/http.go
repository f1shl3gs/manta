package http

import (
	"context"
	"net/http"
	_ "net/http/pprof" // enable http pprof
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/middleware"
	"github.com/f1shl3gs/manta/http/router"
	"github.com/f1shl3gs/manta/multitsdb"
	"github.com/f1shl3gs/manta/raftstore"
	"github.com/f1shl3gs/manta/telemetry/prom"
)

const (
	apiV1Prefix = "/api/v1"
)

type Flusher interface {
	Flush(ctx context.Context) error
}

type Backend struct {
	Flusher Flusher

	router       *router.Router
	PromRegistry *prom.Registry

	BackupService               manta.BackupService
	OrganizationService         manta.OrganizationService
	DashboardService            manta.DashboardService
	UserService                 manta.UserService
	PasswordService             manta.PasswordService
	AuthorizationService        manta.AuthorizationService
	SessionService              manta.SessionService
	OnBoardingService           manta.OnBoardingService
	CheckService                manta.CheckService
	TaskService                 manta.TaskService
	ConfigService               manta.ConfigService
	ScrapeTargetService         manta.ScrapeTargetService
	RegistryService             manta.RegistryService
	NotificationEndpointService manta.NotificationEndpointService
	SecretService               manta.SecretService
	TemplateService             manta.TemplateService
	OperationLogService         manta.OperationLogService

	TenantStorage         multitsdb.TenantStorage
	TenantTargetRetriever multitsdb.TenantTargetRetriever

	ClusterService raftstore.ClusterService
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
	logger = logger.Named("http")

	// assets handler
	assetsHandler, err := NewAssetsHandler(logger)
	if err != nil {
		panic(err)
	}

	// build api handler
	backend.router = router.New(router.Trace(), middleware.Metrics())
	if logger.Core().Enabled(zapcore.DebugLevel) {
		backend.router.Use(middleware.Logging(logger))
	}

	NewOrganizationHandler(backend, logger)
	NewSetupHandler(backend, logger)
	NewSessionHandler(backend.router, logger, backend.UserService, backend.PasswordService, backend.SessionService)
	NewFlushHandler(logger, backend)
	NewDashboardsHandler(backend, logger)
	NewUserHandler(backend, logger)
	NewConfigService(backend, logger)
	NewPromAPIHandler(backend, logger)
	NewScrapeHandler(backend, logger)
	NewRegistryHandler(backend, logger)
	NewChecksHandler(logger, backend.router, backend.CheckService, backend.TaskService, backend.OperationLogService)
	NewTaskHandler(backend, logger)
	NewSecretHandler(logger, backend)
	NewNotificationEendpointHandler(logger, backend)
	NewClusterServiceHandler(logger, backend)
	NewOperationLogHandler(backend)
	NewBuildInfoHandler(backend, logger)

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
	ah.RegisterNoAuthRoute(http.MethodGet, "/")
	ah.RegisterNoAuthRoute(http.MethodGet, "/debug/*wild")
	// TODO: add auth in the future
	ah.RegisterNoAuthRoute(http.MethodGet, configWithID)
	ah.RegisterNoAuthRoute(http.MethodPost, registryPrefix)

	return &Service{
		apiHandler:    ah,
		docHandler:    Redoc(),
		metricHandler: backend.PromRegistry.HTTPHandler(),
		assetsHandler: assetsHandler,
	}
}

func logEncodingError(logger *zap.Logger, r *http.Request, err error) {
	// If we encounter an error while encoding the response to an http request
	// the best thing we can do is logger that error, as we may have already written
	// the headers for the http request in question.
	logger.Info("Error encoding response",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.String("remote", r.RemoteAddr),
		zap.Error(err))
}
