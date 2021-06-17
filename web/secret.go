package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	SecretPrefix  = "/api/v1/secrets"
	SecretKeyPath = "/api/v1/secrets/:key"
)

type SecretHandler struct {
	*Router

	logger        *zap.Logger
	secretService manta.SecretService
}

func NewSecretHandler(logger *zap.Logger, router *Router, secretService manta.SecretService) {
	h := &SecretHandler{
		Router:        router,
		logger:        logger,
		secretService: secretService,
	}

	h.HandlerFunc(http.MethodPut, SecretPrefix, h.handlePut)
	h.HandlerFunc(http.MethodGet, SecretPrefix, h.handleKeys)
	h.HandlerFunc(http.MethodDelete, SecretKeyPath, h.handleDelete)
}

type SecretField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *SecretField) Validate() error {
	if s.Key == "" {
		return errors.New("secret key cannot be empty")
	}

	if s.Value == "" {
		return errors.New("secret value cannot be empty")
	}

	return nil
}

func decodePutSecretRequest(r *http.Request) (string, string, error) {
	var sf SecretField

	err := json.NewDecoder(r.Body).Decode(&sf)
	if err != nil {
		return "", "", err
	}

	return sf.Key, sf.Value, nil
}

func (h *SecretHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	k, v, err := decodePutSecretRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.secretService.PutSecret(ctx, orgID, k, v)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SecretHandler) handleKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	keys, err := h.secretService.GetSecretKeys(ctx, orgID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = encodeResponse(ctx, w, http.StatusOK, keys)
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *SecretHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		return
	}

	key := paramFromRequest(r, "key")
	err = h.secretService.DeleteSecret(ctx, orgID, key)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
