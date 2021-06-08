package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	otclsPrefix = "/api/v1/orgs/:orgID/otcls"
	otclsIDPath = "/api/v1/orgs/:orgID/otcls/:id"
)

type otclHandler struct {
	*Router

	logger      *zap.Logger
	otclService manta.OtclService
}

func otclService(logger *zap.Logger, router *Router, b *Backend) {
	h := &otclHandler{
		Router: router,

		logger:      logger,
		otclService: b.OtclService,
	}

	h.HandlerFunc(http.MethodGet, otclsIDPath, h.getOtcl)
	h.HandlerFunc(http.MethodGet, otclsPrefix, h.getOtcls)
	h.HandlerFunc(http.MethodPost, otclsPrefix, h.createOtcl)
	h.HandlerFunc(http.MethodPatch, otclsIDPath, h.patchOtcl)
	h.HandlerFunc(http.MethodDelete, otclsIDPath, h.deleteOtcl)
}

func (h *otclHandler) getOtcls(w http.ResponseWriter, r *http.Request) {
	var (
		orgID manta.ID
		ctx   = r.Context()
	)

	if err := orgID.DecodeFromString(r.URL.Query().Get("orgID")); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	otcls, err := h.otclService.FindOtcls(ctx, manta.OtclFilter{OrgID: &orgID})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, otcls); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *otclHandler) getOtcl(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	c, err := h.otclService.FindOtclByID(r.Context(), id)
	if err != nil {
		h.HandleHTTPError(r.Context(), err, w)
		return
	}

	accept := r.Header.Get("Accept")
	switch {
	case strings.Contains(accept, "application/json"):
		err = encodeResponse(ctx, w, http.StatusOK, c)
	case strings.Contains(accept, "application/octet-stream"):
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(c.Content))
	default:
		// On chrome "application/yaml" will download, while "text/yaml" will display
		w.Header().Set("Content-Type", "text/yaml")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(c.Content))
	}

	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *otclHandler) decodeOtclRequest(r *http.Request) (*manta.Otcl, error) {
	otcl := &manta.Otcl{}
	err := json.NewDecoder(r.Body).Decode(otcl)
	if err != nil {
		return nil, err
	}

	if !otcl.OrgID.Valid() {
		return nil, manta.ErrInvalidOrgID
	}

	return otcl, nil
}

func (h *otclHandler) createOtcl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	otcl, err := h.decodeOtclRequest(r)
	if err != nil {
		// todo: handle bad request
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.otclService.CreateOtcl(ctx, otcl)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusCreated, otcl); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeOtclPatch(r *http.Request, p *manta.OtclPatch) error {
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		return err
	}

	if p.Name == nil && p.Content == nil && p.Desc == nil {
		return &manta.Error{
			Code: manta.EInvalid,
			Msg:  "Patch is empty",
			Op:   "validate",
		}
	}

	return nil
}

// patchOtcl updates a Otcl
func (h *otclHandler) patchOtcl(w http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		patch manta.OtclPatch
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = decodeOtclPatch(r, &patch)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	otcl, err := h.otclService.PatchOtcl(ctx, id, patch)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, otcl); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *otclHandler) deleteOtcl(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.otclService.DeleteOtcl(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
