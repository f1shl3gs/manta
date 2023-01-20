package http

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"
	"github.com/f1shl3gs/manta/pkg/tracing"
)

type AuthenticationHandler struct {
	logger *zap.Logger

	AuthorizationService manta.AuthorizationService
	UserService          manta.UserService
	SessionService       manta.SessionService

	noAuthRouter *httprouter.Router
	handler      http.Handler
	errorHandler manta.HTTPErrorHandler
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	span, ctx := tracing.StartSpanFromContext(r.Context())
	defer span.Finish()

	var (
		a    manta.Authorizer
		err  error
		path = r.URL.Path
	)

	if handler, _, _ := h.noAuthRouter.Lookup(r.Method, path); handler != nil {
		h.handler.ServeHTTP(w, r)
		return
	}

	switch probeAuthType(r) {
	case "token":
		// TODO: implement
	case "session":
		a, err = h.extractSession(ctx, r)
	default:
		h.handleUnauthorized(w, r, &manta.Error{
			Code: manta.EUnauthorized,
			Msg:  "no auth type found",
		})
		return
	}

	if err == manta.ErrSessionExpired || err == manta.ErrSessionNotFound {
		h.handleUnauthorized(w, r, err)
		return
	}

	if err != nil {
		h.errorHandler.HandleHTTPError(ctx, err, w)
		return
	}

	span.SetTag("user_id", a.GetUserID().String())

	ctx = authorizer.SetAuthorizer(ctx, a)
	r = r.WithContext(ctx)
	h.handler.ServeHTTP(w, r)
}

func (h *AuthenticationHandler) RegisterNoAuthRoute(method, path string) {
	h.noAuthRouter.HandlerFunc(method, path, func(w http.ResponseWriter, r *http.Request) {})
}

func probeAuthType(r *http.Request) string {
	if v := r.Header.Get("Authorization"); v != "" {
		return "token"
	}

	for _, c := range r.Cookies() {
		if c.Name == SessionCookieKey {
			return "session"
		}
	}

	if r.URL.Query().Has("token") {
		return "token"
	}

	return ""
}

func (h *AuthenticationHandler) extractSession(ctx context.Context, r *http.Request) (manta.Authorizer, error) {
	c, err := r.Cookie(SessionCookieKey)
	if err != nil {
		return nil, err
	}

	var id manta.ID
	err = id.DecodeFromString(c.Value)
	if err != nil {
		return nil, err
	}

	session, err := h.SessionService.FindSession(ctx, id)
	if err != nil {
		return nil, err
	}

	err = h.SessionService.RenewSession(ctx, id, time.Now().Add(12*time.Hour))
	if err != nil {
		h.logger.Error("renew session failed",
			zap.Error(err),
			zap.String("userId", id.String()))
	}

	return session, nil
}

func (h *AuthenticationHandler) handleUnauthorized(w http.ResponseWriter, r *http.Request, err error) {
	h.logger.Debug("Unauthorized http request",
		zap.String("remote", r.RemoteAddr),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Error(err))

	w.WriteHeader(http.StatusUnauthorized)
}
