package http

import (
	"encoding/json"
	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
	"net/http"
)

const setupPath = apiV1Prefix + "/setup"

type SetupHandler struct {
	*Router

	logger            *zap.Logger
	onBoardingService manta.OnBoardingService
}

func NewSetupHandler(backend *Backend, logger *zap.Logger) {
	h := &SetupHandler{
		Router:            backend.router,
		logger:            logger,
		onBoardingService: backend.OnBoardingService,
	}

	h.HandlerFunc(http.MethodPost, setupPath, h.onBoarding)
}

func (h *SetupHandler) onBoarding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	br, err := decodeOnBoardingRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	result, err := h.onBoardingService.Setup(ctx, br)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	h.logger.Info("OnBoarding successfully",
		zap.String("org", br.Organization),
		zap.String("username", br.Username),
	)

	if err = encodeResponse(ctx, w, http.StatusOK, result); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeOnBoardingRequest(r *http.Request) (*manta.OnBoardingRequest, error) {
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
