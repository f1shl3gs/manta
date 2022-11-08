package http

import (
	"context"
	"net/http"
    "strings"

    "github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/middlewares"
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
    ConfigurationService manta.ConfigurationService
}

type Service struct {
	apiHandler http.Handler

    assetsHandler http.Handler
}

func (s *Service ) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path

    if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/debug") {
        s.apiHandler.ServeHTTP(w, r)
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

	ah := &AuthenticationHandler{
		logger:               logger,
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
	ah.RegisterNoAuthRoute(http.MethodGet, DebugFlushPath)
    ah.RegisterNoAuthRoute(http.MethodGet, "/")

	// set kinds of global middleware
	handler := http.Handler(ah)
    handler = middlewares.Trace(handler)

	// enable access log middleware
	if logger.Core().Enabled(zapcore.DebugLevel) {
		handler = middlewares.Logging(logger, handler)
	}

	return &Service{
		apiHandler: handler,
        assetsHandler: assetsHandler,
	}
}
