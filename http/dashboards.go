package http

import (
	"encoding/json"
	"net/http"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

const (
	dashboardsPrefix     = apiV1Prefix + `/dashboards`
	dashboardsWithID     = dashboardsPrefix + `/:id`
	dashboardCellsPrefix = dashboardsWithID + `/cells`
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
		logger:              logger.With(zap.String("handler", "dashboard")),
		organizationService: backend.OrganizationService,
		dashboardService:    backend.DashboardService,
	}

	h.HandlerFunc(http.MethodGet, dashboardsPrefix, h.listDashboard)
	h.HandlerFunc(http.MethodGet, dashboardsWithID, h.getDashboard)
	h.HandlerFunc(http.MethodPost, dashboardsPrefix, h.create)
	h.HandlerFunc(http.MethodDelete, dashboardsWithID, h.delete)
	h.HandlerFunc(http.MethodPatch, dashboardsWithID, h.updateMeta)

	h.HandlerFunc(http.MethodPost, dashboardCellsPrefix, h.createCell)

	return h
}

func (h *DashboardsHandler) listDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIdFromQuery(r)
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

func (h *DashboardsHandler) getDashboard(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	d, err := h.dashboardService.FindDashboardByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, d); err != nil {
		logEncodingError(h.logger, r, err)
	}
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

	id, err := idFromPath(r)
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

	id, err := idFromPath(r)
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

func decodeCreateCell(r *http.Request) (*manta.Cell, error) {
	var cc struct {
		W, H, X, Y int32
	}

	err := json.NewDecoder(r.Body).Decode(&cc)
	if err != nil {
		return nil, err
	}

	return &manta.Cell{
		W: cc.W,
		H: cc.H,
		X: cc.X,
		Y: cc.Y,
	}, nil
}

func (h *DashboardsHandler) createCell(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cell, err := decodeCreateCell(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := h.dashboardService.AddDashboardCell(ctx, id, cell); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
