package web

import (
	"encoding/json"
	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
	"net/http"
)

const (
	signinPath  = "/api/v1/signin"
	signoutPath = "/api/v1/signout"
)

type SessionHandler struct {
	*Router

	logger          *zap.Logger
	userService     manta.UserService
	passwordService manta.PasswordService
}

func NewSessionHandler(logger *zap.Logger, userService manta.UserService) *SessionHandler {
	h := &SessionHandler{
		logger: logger,
	}

	h.HandlerFunc(http.MethodPost, signinPath, h.handleSignin)

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
		h.HandleHTTPError(ctx, err, w)
		return
	}
}
