package web

import (
	"encoding/json"
	"github.com/f1shl3gs/manta"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
)

const (
	DashboardPrefix     = "/api/v1/dashboards"
	DashboardIDPath     = "/api/v1/dashboards/:id"
	DashboardCellPrefix = "/api/v1/dashboards/:id/cells"
	DashboardCellIDPath = "/api/v1/dashboards/:id/cells/:cid"
)

type DashboardHandler struct {
	*Router

	logger           *zap.Logger
	dashboardService manta.DashboardService
}

func NewDashboardService(h *DashboardHandler) {
	h.HandlerFunc(http.MethodGet, DashboardPrefix, h.list)
	h.HandlerFunc(http.MethodGet, DashboardIDPath, h.get)
	h.HandlerFunc(http.MethodPost, DashboardPrefix, h.handleCreate)
	h.HandlerFunc(http.MethodPost, DashboardCellPrefix, h.handleAddCell)
}

func (h *DashboardHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ds, err := h.dashboardService.FindDashboards(ctx, manta.DashboardFilter{})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, &ds); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromRequestPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	d, err := h.dashboardService.FindDashboardByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, d); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromRequestPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.dashboardService.DeleteDashboard(ctx, id); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusNoContent, nil); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var (
		dash manta.Dashboard
		ctx  = r.Context()
	)

	if err := json.NewDecoder(r.Body).Decode(&dash); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err := h.dashboardService.CreateDashboard(ctx, &dash)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusCreated, &dash); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardHandler) handleAddCell(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		params = httprouter.ParamsFromContext(ctx)
	)

	id, err := idFromParams(params, "id")
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
