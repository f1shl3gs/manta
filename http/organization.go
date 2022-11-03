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
	h.HandlerFunc(http.MethodGet, organizationPrefix+"/:orgId", h.getOrganization)

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

func (h *OrganizationHandler) getOrganization(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		id  manta.ID
	)

	if err := id.DecodeFromString(ExtractParamFromContext(ctx, "orgId")); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	org, err := h.organizationService.FindOrganizationByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, org); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *OrganizationHandler) deleteOrganization(w http.ResponseWriter, r *http.Request) {

}