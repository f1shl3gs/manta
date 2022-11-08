package http

import (
	"net/http"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

const (
	UserPrefix = apiV1Prefix + "/users"
)

type UsersHandler struct {
	*Router
	logger *zap.Logger

	userService manta.UserService
}

func NewUserHandler(backend *Backend, logger *zap.Logger) *UsersHandler {
	h := &UsersHandler{
		Router:      backend.router,
		logger:      logger.With(zap.String("handler", "users")),
		userService: backend.UserService,
	}

	h.HandlerFunc(http.MethodGet, UserPrefix, h.list)

	return h
}

func (h *UsersHandler) list(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	users, err := h.userService.FindUsers(ctx, manta.UserFilter{})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, users); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
