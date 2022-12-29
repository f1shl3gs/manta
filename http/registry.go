package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"
)

const (
	registryPrefix = apiV1Prefix + "/registry"
)

type RegistryHandler struct {
	*router.Router

	logger          *zap.Logger
	registryService manta.RegistryService
}

func NewRegistryService(backend *Backend, logger *zap.Logger) {
	h := &RegistryHandler{
		Router:          backend.router,
		logger:          logger.With(zap.String("handler", "registry")),
		registryService: backend.RegistryService,
	}

	h.HandlerFunc(http.MethodPost, registryPrefix, h.register)
	h.HandlerFunc(http.MethodGet, registryPrefix, h.catalog)
}

func (h *RegistryHandler) register(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		ins = &manta.Instance{}
	)

	err := json.NewDecoder(r.Body).Decode(&ins)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.registryService.Register(ctx, ins)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RegistryHandler) catalog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := h.registryService.Catalog(ctx)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, &list); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
