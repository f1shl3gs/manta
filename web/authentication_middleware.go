package web

import (
	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

type AuthenticationHandler struct {
	logger *zap.Logger

	AuthorizationService manta.AuthorizationService
	UserService          manta.UserService
}
