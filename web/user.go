package web

import (
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"
	"go.uber.org/zap"
)

const (
	userPrefix = "/api/v1/users"
	viewerPath = "/api/v1/viewer"
)

type UserHandler struct {
	*Router

	logger      *zap.Logger
	userService manta.UserService
}

func (h *UserHandler) viewerHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	auth := authorizer.FromContext(r.Context())
	if auth == nil {
		encodeResponse(ctx, w, http.StatusUnauthorized, nil)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, auth); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
