package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"
)

var (
	secretPrefix  = apiV1Prefix + "/secrets"
	secretKeyPath = secretPrefix + "/:key"
)

type SecretHandler struct {
	*router.Router
	logger *zap.Logger

	secretService manta.SecretService
}

func NewSecretHandler(logger *zap.Logger, backend *Backend) {
	h := &SecretHandler{
		Router:        backend.router,
		logger:        logger.With(zap.String("handler", "secret")),
		secretService: backend.SecretService,
	}

	h.HandlerFunc(http.MethodGet, secretPrefix, h.handleList)
	h.HandlerFunc(http.MethodPost, secretPrefix, h.handlePut)
	h.HandlerFunc(http.MethodDelete, secretKeyPath, h.handleDelete)
}

func (h *SecretHandler) handleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIDFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	secrets, err := h.secretService.GetSecrets(ctx, orgID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, secrets); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeSecret(r *http.Request) (*manta.Secret, error) {
	secret := &manta.Secret{}

	err := json.NewDecoder(r.Body).Decode(secret)
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "unmarshal secret failed",
			Err:  err,
		}
	}

	return secret, nil
}

func (h *SecretHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	secret, err := decodeSecret(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	secret, err = h.secretService.PutSecret(ctx, secret)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusCreated, secret); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *SecretHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIDFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	params := httprouter.ParamsFromContext(ctx)
	key := params.ByName("key")

	err = h.secretService.DeleteSecret(ctx, orgID, key)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}
