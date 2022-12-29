package http

import (
	"encoding/json"
    "net/http"

	"go.uber.org/zap"

    "github.com/f1shl3gs/manta/http/router"
	"github.com/f1shl3gs/manta"
)

const setupPath = apiV1Prefix + "/setup"

type SetupHandler struct {
	*router.Router

	logger            *zap.Logger
	onBoardingService manta.OnBoardingService
	sessionService    manta.SessionService
}

func NewSetupHandler(backend *Backend, logger *zap.Logger) {
	h := &SetupHandler{
		Router:            backend.router,
		logger:            logger.With(zap.String("handler", "setup")),
		onBoardingService: backend.OnBoardingService,
		sessionService:    backend.SessionService,
	}

	h.HandlerFunc(http.MethodGet, setupPath, h.onBoarded)
	h.HandlerFunc(http.MethodPost, setupPath, h.onBoarding)
}

type OnboardedResult struct {
	Allow bool `json:"allow"`
}

func (h *SetupHandler) onBoarded(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	allow, err := h.onBoardingService.Onboarded(ctx)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, OnboardedResult{Allow: !allow}); err != nil {
		logEncodingError(h.logger, r, err)
	}
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

	if sess, err := h.sessionService.CreateSession(ctx, result.User.ID); err != nil {
		h.logger.Error("Create initial session failed", zap.Error(err))
		h.HandleHTTPError(ctx, err, w)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     SessionCookieKey,
			Value:    sess.ID.String(),
			HttpOnly: true,
			// Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})

		if err = h.EncodeResponse(ctx, w, http.StatusOK, result); err != nil {
			logEncodingError(h.logger, r, err)
		}
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
