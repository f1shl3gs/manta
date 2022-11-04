package http

import (
	"encoding/json"
	"net/http"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

const (
	dashboardsPrefix = apiV1Prefix + `/dashboards`
    dashboardsWithID = dashboardsPrefix + `/:id`
)

type DashboardsHandler struct {
	*Router

	logger              *zap.Logger
	organizationService manta.OrganizationService
	dashboardService    manta.DashboardService
}

func NewDashboardsHandler(backend *Backend, logger *zap.Logger) *DashboardsHandler {
	h := &DashboardsHandler{
		Router:              backend.router,
		logger:              logger,
		organizationService: backend.OrganizationService,
		dashboardService:    backend.DashboardService,
	}

	h.HandlerFunc(http.MethodGet, dashboardsPrefix, h.list)
	h.HandlerFunc(http.MethodPost, dashboardsPrefix, h.create)
    h.HandlerFunc(http.MethodDelete, dashboardsWithID, h.delete)
    h.HandlerFunc(http.MethodPatch, dashboardsWithID, h.updateMeta)

	return h
}

func (h *DashboardsHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := OrgIdFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	dashboards, err := h.dashboardService.FindDashboards(ctx, manta.DashboardFilter{OrganizationID: &orgID})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, dashboards); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeDashboard(r *http.Request) (*manta.Dashboard, error) {
	var dashboard manta.Dashboard

	err := json.NewDecoder(r.Body).Decode(&dashboard)
	if err != nil {
		return nil, err
	}

	// force reset orgID!?

	return &dashboard, nil
}

func (h *DashboardsHandler) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dashboard, err := decodeDashboard(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := h.dashboardService.CreateDashboard(ctx, dashboard); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusCreated, dashboard); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardsHandler) delete(w http.ResponseWriter, r *http.Request) {
    var ctx = r.Context()

    id, err := IDFromPath(r)
    if err != nil {
        h.HandleHTTPError(ctx, err, w)
        return
    }

    err = h.dashboardService.DeleteDashboard(ctx, id)
    if err != nil {
        h.HandleHTTPError(ctx, err, w)
        return
    }
}

func (h *DashboardsHandler) updateMeta(w http.ResponseWriter, r *http.Request) {
    var (
        ctx = r.Context()
        upd manta.DashboardUpdate
    )

    id, err := IDFromPath(r)
    if err != nil {
        h.HandleHTTPError(ctx, err, w)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
        h.HandleHTTPError(ctx, err, w)
        return
    }

    if _, err = h.dashboardService.UpdateDashboard(ctx, id, upd); err != nil {
        h.HandleHTTPError(ctx, err, w)
    }
}
