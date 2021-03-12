package web

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type HealthzHandler struct {
	logger *zap.Logger
}

func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
	})

	if err != nil {
		h.logger.Error("Write healthz response failed",
			zap.Error(err))
	}
}

func newHealthzHandler(logger *zap.Logger) http.Handler {
	return &HealthzHandler{logger: logger}
}
