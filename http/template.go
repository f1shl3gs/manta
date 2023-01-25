package http

import (
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"

	"go.uber.org/zap"
)

const (
	templatePrefix = apiV1Prefix + "/templates"
	templateIDPath = templatePrefix + "/:id"
)

type TemplateHandler struct {
	*router.Router
	logger *zap.Logger

	templateService manta.TemplateService
}

func NewTemplateHandler(backend *Backend, logger *zap.Logger) {
	h := &TemplateHandler{
		Router: backend.router,
		logger: logger,

		templateService: backend.TemplateService,
	}

	h.HandlerFunc(http.MethodGet, templatePrefix, h.list)
}

func (h *TemplateHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIdFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	templates, err := h.templateService.ListTemplate(ctx, orgID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := h.EncodeResponse(ctx, w, http.StatusOK, templates); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
