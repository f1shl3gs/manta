package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorization"
	"github.com/f1shl3gs/manta/pkg/tracing"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const tokenScheme = "Bearer "

var (
	ErrAuthHeaderMissing = errors.New("authorization Header is missing")
	ErrAuthBadScheme     = errors.New("authorization Header Scheme is invalid")
)

type AuthenticationHandler struct {
	logger *zap.Logger

	AuthorizationService manta.AuthorizationService
	UserService          manta.UserService
	Keyring              manta.Keyring

	noAuthRouter *httprouter.Router
	handler      http.Handler
	errorHandler manta.HTTPErrorHandler
	tokenParser  *authorization.TokenParser
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, _, _ := h.noAuthRouter.Lookup(r.Method, r.URL.Path); handler != nil {
		h.handler.ServeHTTP(w, r)
		return
	}

	var (
		ctx        = r.Context()
		authorizer manta.Authorizer
		err        error
	)

	switch probeAuthType(r) {
	case "token":
	case "jwt":
		authorizer, err = h.extractJWT(ctx, r)
	default:
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		h.handleUnauthorized(ctx, w, err)
		return
	}

	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	span.SetTag("user_id", authorizer.GetUserID().String())

	h.handler.ServeHTTP(w, r)
}

func probeAuthType(r *http.Request) string {
	return "jwt"
}

func (h *AuthenticationHandler) extractJWT(ctx context.Context, r *http.Request) (manta.Authorizer, error) {
	c, err := r.Cookie("jwt")
	if err != nil {
		return nil, err
	}

	t, err := h.tokenParser.Parse(c.Value)
	if err != nil {
		return nil, err
	}

	return h.AuthorizationService.FindAuthorizationByID(ctx, t.Identifier())
}

func (h *AuthenticationHandler) handleUnauthorized(ctx context.Context, w http.ResponseWriter, err error) {
	h.logger.Info("unauthorized", zap.Error(err))
	h.errorHandler.HandleHTTPError(ctx, err, w)
}
