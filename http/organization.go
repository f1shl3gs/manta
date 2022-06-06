package http

import (
	"net/http"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

const (
	organizationPrefix = apiV1Prefix + "/organizations"
)

type OrganizationHandler struct {
	*Router

	logger              *zap.Logger
	organizationService manta.OrganizationService
}

func NewOrganizationHandler(backend *Backend, logger *zap.Logger) *OrganizationHandler {
	h := &OrganizationHandler{
		Router:              backend.router,
		logger:              logger.With(zap.String("handler", "organization")),
		organizationService: backend.OrganizationService,
	}

	h.HandlerFunc(http.MethodGet, organizationPrefix, h.listOrganizations)

	return h
}

func (h *OrganizationHandler) listOrganizations(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	organizations, _, err := h.organizationService.FindOrganizations(ctx, manta.OrganizationFilter{})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, organizations); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
