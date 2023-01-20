package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"
)

const (
	UserPrefix = apiV1Prefix + "/users"
)

type UsersHandler struct {
	*router.Router
	logger *zap.Logger

	userService    manta.UserService
	sessionService manta.SessionService
}

func NewUserHandler(backend *Backend, logger *zap.Logger) *UsersHandler {
	h := &UsersHandler{
		Router:         backend.router,
		logger:         logger.With(zap.String("handler", "users")),
		userService:    backend.UserService,
		sessionService: backend.SessionService,
	}

	h.HandlerFunc(http.MethodGet, UserPrefix, h.list)

	return h
}

func (h *UsersHandler) list(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	type UserExt struct {
		*manta.User
		LastSeen *time.Time `json:"lastSeen,omitempty"`
	}

	users, err := h.userService.FindUsers(ctx, manta.UserFilter{})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	resps := make([]UserExt, len(users))
	// TODO: implement something like GetBatchXXXX to improve performance ?
	for i, u := range users {
		resps[i].User = u
		session, err := h.sessionService.FindSession(ctx, u.ID)
		if err == nil {
			resps[i].LastSeen = &session.LastSeen
		}
	}

	if err := h.EncodeResponse(ctx, w, http.StatusOK, resps); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
