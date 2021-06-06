package web

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	orgsPrefix = "/api/v1/orgs"
	orgsIDPath = "/api/v1/orgs/:id"
)

type OrganizationHandler struct {
	*Router

	logger              *zap.Logger
	OrganizationService manta.OrganizationService
}

func NewOrganizationHandler(logger *zap.Logger, router *Router, b *Backend) {
	h := &OrganizationHandler{
		Router:              router,
		logger:              logger,
		OrganizationService: b.OrganizationService,
	}

	h.HandlerFunc(http.MethodGet, orgsPrefix, h.handleGetOrgs)
	h.HandlerFunc(http.MethodPost, orgsPrefix, h.handleCreateOrg)
	h.HandlerFunc(http.MethodDelete, orgsIDPath, h.handleDeleteOrg)
	h.HandlerFunc(http.MethodGet, orgsIDPath, h.handleGetOrg)
}

func (h *OrganizationHandler) handleGetOrgs(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	orgs, _, err := h.OrganizationService.FindOrganizations(ctx, manta.OrganizationFilter{})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, orgs); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeOrganization(r io.Reader) (*manta.Organization, error) {
	org := &manta.Organization{}

	dec := json.NewDecoder(r)
	if err := dec.Decode(org); err != nil {
		return nil, err
	}

	if err := org.Validate(); err != nil {
		return nil, err
	}

	return org, nil
}

func (h *OrganizationHandler) handleCreateOrg(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	org, err := decodeOrganization(r.Body)
	if err != nil {
		err = &manta.Error{
			Code: manta.EInvalid,
			Msg:  "decode organization request failed",
			Err:  err,
		}

		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.OrganizationService.CreateOrganization(ctx, org)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, org); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *OrganizationHandler) handleDeleteOrg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.OrganizationService.DeleteOrganization(ctx, id); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OrganizationHandler) handleGetOrg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	org, err := h.OrganizationService.FindOrganizationByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, org); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
