package web

import (
	"io"
	"net/http"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

const (
	keyringPrefix = "/api/v1/keyring"
)

type keyringHandler struct {
	*Router

	logger  *zap.Logger
	keyring manta.Keyring
}

func NewKeyringHandler(router *Router, logger *zap.Logger, keyring manta.Keyring) {
	h := &keyringHandler{
		Router:  router,
		logger:  logger,
		keyring: keyring,
	}

	h.HandlerFunc(http.MethodPut, keyringPrefix, h.handleAdd)
}

func (h *keyringHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = h.keyring.AddKey(r.Context(), data)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
