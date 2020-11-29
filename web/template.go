package web

import (
	"net/http"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

type TemplateHandler struct {
	logger *zap.Logger

	CollectionService           manta.OtclService
	DashboardService            manta.DashboardService
	CheckService                manta.CheckService
	NotificationEndpointService manta.NotificationEndpointService
}

func NewTemplateHandler(logger *zap.Logger, router *Router) {
	h := &TemplateHandler{
		logger: logger.With(zap.String("handler", "template")),
	}

	router.HandlerFunc(http.MethodPost, "/api/v1/templates/apply", h.apply)
}

func (h *TemplateHandler) apply(w http.ResponseWriter, r *http.Request) {

}
