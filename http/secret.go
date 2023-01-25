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
	secretPrefix      = apiV1Prefix + "/secrets"
	secretKeyPath     = secretPrefix + "/:key"
	secretChangesPath = secretKeyPath + "/changes"
)

type SecretHandler struct {
	*router.Router
	logger *zap.Logger

	secretService manta.SecretService
	oplogService  manta.OperationLogService
}

func NewSecretHandler(logger *zap.Logger, backend *Backend) {
	h := &SecretHandler{
		Router:        backend.router,
		logger:        logger.With(zap.String("handler", "secret")),
		secretService: backend.SecretService,
		oplogService:  backend.OperationLogService,
	}

	h.HandlerFunc(http.MethodGet, secretPrefix, h.handleList)
	h.HandlerFunc(http.MethodPost, secretPrefix, h.handlePut)
	h.HandlerFunc(http.MethodDelete, secretKeyPath, h.handleDelete)
	h.HandlerFunc(http.MethodGet, secretChangesPath, h.handleListChanges)
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

func extractSecretKey(r *http.Request) string {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName("key")
}

func (h *SecretHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIDFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	key := extractSecretKey(r)

	err = h.secretService.DeleteSecret(ctx, orgID, key)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}

func (h *SecretHandler) handleListChanges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := extractSecretKey(r)

	opts, err := manta.DecodeFindOptions(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	changes, _, err := h.oplogService.FindOperationLogsByID(ctx, manta.UniqueKeyToID(key), opts)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, changes); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
