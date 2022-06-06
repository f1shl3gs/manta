package http

import (
	"context"
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/middlewares"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
}

type Service struct {
	http.Handler

	backend *Backend
}

func New(logger *zap.Logger, backend *Backend) *Service {
	backend.router = &Router{
		Router: httprouter.New(),
		logger: logger,
	}

	NewOrganizationHandler(backend, logger)
	NewSetupHandler(backend, logger)
	NewSessionHandler(backend.router, logger, backend.UserService, backend.PasswordService, backend.SessionService)
	NewFlushHandler(logger, backend.router, backend.Flusher)

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
	ah.RegisterNoAuthRoute(http.MethodPost, signinPath)

	// set kinds of global middleware
	handler := http.Handler(ah)

	// enable access log middleware
	if logger.Core().Enabled(zapcore.DebugLevel) {
		handler = middlewares.Logging(logger, handler)
	}

	return &Service{
		Handler: handler,
		backend: backend,
	}
}
