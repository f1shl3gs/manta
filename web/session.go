package web

import (
	"context"
	"encoding/json"
	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
	"net/http"
)

const (
	signinPath  = "/api/v1/signin"
	signoutPath = "/api/v1/signout"

	SessionCookieKey = "manta_session"
)

type SessionHandler struct {
	*Router

	logger          *zap.Logger
	userService     manta.UserService
	passwordService manta.PasswordService
	sessionService  manta.SessionService
}

func NewSessionHandler(
	router *Router,
	logger *zap.Logger,
	userService manta.UserService,
	passwordService manta.PasswordService,
	sessionService manta.SessionService,
) *SessionHandler {
	h := &SessionHandler{
		Router:          router,
		logger:          logger,
		userService:     userService,
		passwordService: passwordService,
		sessionService:  sessionService,
	}

	h.HandlerFunc(http.MethodPost, signinPath, h.handleSignin)
	h.HandlerFunc(http.MethodDelete, signoutPath, h.handleSignout)

	return h
}

type signinReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func decodeSigninRequest(r *http.Request) (*signinReq, error) {
	s := &signinReq{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// handleSignin is the HTTP handler for the POST /api/v1/signin route
func (h *SessionHandler) handleSignin(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	sr, err := decodeSigninRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	u, err := h.userService.FindUser(ctx, manta.UserFilter{
		Name: &sr.Username,
	})

	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.passwordService.ComparePassword(ctx, u.ID, sr.Password)
	if err != nil {
		if err == manta.ErrPasswordNotMatch {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.HandleHTTPError(ctx, err, w)
		return
	}

	session, err := h.sessionService.CreateSession(ctx, u.ID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.encodeCookie(ctx, w, session.ID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
	}
}

func (h *SessionHandler) handleSignout(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		id  manta.ID
	)

	c, err := r.Cookie(SessionCookieKey)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = id.DecodeFromString(c.Value)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.sessionService.RevokeSession(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SessionHandler) encodeCookie(ctx context.Context, w http.ResponseWriter, id manta.ID) error {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieKey,
		Value:    id.String(),
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}
