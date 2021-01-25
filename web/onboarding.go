package web

import (
	"encoding/json"
	"net/http"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

const (
	setupPath = "/api/v1/setup"
)

type SetupHandler struct {
	*Router

	logger            *zap.Logger
	onBoardingService manta.OnBoardingService
}

func NewSetupHandler(router *Router, logger *zap.Logger, backend *Backend) {
	h := &SetupHandler{
		Router: router,
		logger: logger,
		onBoardingService: manta.NewOnBoardingService(
			backend.UserService,
			backend.PasswordService,
			backend.AuthorizationService,
			backend.OrganizationService),
	}

	h.HandlerFunc(http.MethodPost, setupPath, h.onBoarding)
}

// onBoarding is the HTTP handler for /api/v1/setup
func (h *SetupHandler) onBoarding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	br, err := h.decodeOnBoardingRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.onBoardingService.Setup(ctx, br)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}

func (h *SetupHandler) decodeOnBoardingRequest(r *http.Request) (*manta.OnBoardingRequest, error) {
	br := &manta.OnBoardingRequest{}

	err := json.NewDecoder(r.Body).Decode(br)
	if err != nil {
		return nil, err
	}

	err = br.Validate()
	if err != nil {
		return nil, err
	}

	return br, nil
}
