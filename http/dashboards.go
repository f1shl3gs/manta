package http

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

const (
	dashboardsPrefix     = apiV1Prefix + `/dashboards`
	dashboardsWithID     = dashboardsPrefix + `/:id`
	dashboardCellsPrefix = dashboardsWithID + `/cells`
	dashboardCellIDPath  = dashboardCellsPrefix + `/:cellId`
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

	h.HandlerFunc(http.MethodGet, dashboardCellIDPath, h.getCell)
	h.HandlerFunc(http.MethodPost, dashboardCellsPrefix, h.addCell)
	h.HandlerFunc(http.MethodPatch, dashboardCellIDPath, h.updateCell)
	h.HandlerFunc(http.MethodPut, dashboardCellsPrefix, h.replaceCells)
	h.HandlerFunc(http.MethodDelete, dashboardCellIDPath, h.deleteCell)

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

func (h *DashboardsHandler) addCell(w http.ResponseWriter, r *http.Request) {
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

func (h *DashboardsHandler) getCell(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	dashId, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cellId, err := cellIdFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cell, err := h.dashboardService.GetDashboardCell(ctx, dashId, cellId)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, cell); err != nil {
		h.HandleHTTPError(ctx, err, w)
	}
}

func (h *DashboardsHandler) updateCell(w http.ResponseWriter, r *http.Request) {
	var (
		upd manta.DashboardCellUpdate
		ctx = r.Context()
	)

	dashboardId, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cellId, err := cellIdFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	_, err = h.dashboardService.UpdateDashboardCell(ctx, dashboardId, cellId, upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func cellIdFromPath(r *http.Request) (manta.ID, error) {
	var id manta.ID

	text := extractParamFromContext(r.Context(), "cellId")

	return id, id.DecodeFromString(text)
}

func decodeCells(r *http.Request) ([]manta.Cell, error) {
	var (
		cells []manta.Cell
		err   error
	)

	err = json.NewDecoder(r.Body).Decode(&cells)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(cells); i++ {
		if err = cells[i].Validate(); err != nil {
			return nil, errors.Wrapf(err, "invalid cell at %d", i)
		}
	}

	return cells, err
}

func (h *DashboardsHandler) replaceCells(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cells, err := decodeCells(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.dashboardService.ReplaceDashboardCells(ctx, id, cells)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DashboardsHandler) deleteCell(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	did, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cid, err := cellIdFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.dashboardService.RemoveDashboardCell(ctx, did, cid); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}
