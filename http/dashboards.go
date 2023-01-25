package http

import (
	"encoding/json"
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	dashboardsPrefix     = apiV1Prefix + `/dashboards`
	dashboardsWithID     = dashboardsPrefix + `/:id`
	dashboardChangesPath = dashboardsWithID + `/changes`
	dashboardCellsPrefix = dashboardsWithID + `/cells`
	dashboardCellIDPath  = dashboardCellsPrefix + `/:cellId`
)

type DashboardsHandler struct {
	*router.Router

	logger              *zap.Logger
	organizationService manta.OrganizationService
	dashboardService    manta.DashboardService
	operationLogService manta.OperationLogService
}

func NewDashboardsHandler(backend *Backend, logger *zap.Logger) *DashboardsHandler {
	h := &DashboardsHandler{
		Router:              backend.router,
		logger:              logger.With(zap.String("handler", "dashboard")),
		organizationService: backend.OrganizationService,
		dashboardService:    backend.DashboardService,
		operationLogService: backend.OperationLogService,
	}

	h.HandlerFunc(http.MethodGet, dashboardsPrefix, h.listDashboard)
	h.HandlerFunc(http.MethodGet, dashboardsWithID, h.getDashboard)
	h.HandlerFunc(http.MethodPost, dashboardsPrefix, h.createDashboard)
	h.HandlerFunc(http.MethodDelete, dashboardsWithID, h.deletedashboard)
	h.HandlerFunc(http.MethodPatch, dashboardsWithID, h.updateMeta)

	h.HandlerFunc(http.MethodGet, dashboardCellIDPath, h.getCell)
	h.HandlerFunc(http.MethodPost, dashboardCellsPrefix, h.addCell)
	h.HandlerFunc(http.MethodPatch, dashboardCellIDPath, h.updateCell)
	h.HandlerFunc(http.MethodPut, dashboardCellsPrefix, h.replaceCells)
	h.HandlerFunc(http.MethodDelete, dashboardCellIDPath, h.deleteCell)

	h.HandlerFunc(http.MethodGet, dashboardChangesPath, h.changes)

	return h
}

func (h *DashboardsHandler) listDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIDFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	dashboards, err := h.dashboardService.FindDashboards(ctx, manta.DashboardFilter{OrgID: orgID})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := h.EncodeResponse(ctx, w, http.StatusOK, dashboards); err != nil {
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

	if err := h.EncodeResponse(ctx, w, http.StatusOK, d); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardsHandler) createDashboard(w http.ResponseWriter, r *http.Request) {
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

	if err := h.EncodeResponse(ctx, w, http.StatusCreated, dashboard); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardsHandler) deletedashboard(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.dashboardService.RemoveDashboard(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeDashboardUpdate(r *http.Request) (manta.DashboardUpdate, error) {
	var upd manta.DashboardUpdate

	err := json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		return manta.DashboardUpdate{}, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "invalid dashboard update",
			Err:  err,
		}
	}

	return upd, nil
}

func (h *DashboardsHandler) updateMeta(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	upd, err := decodeDashboardUpdate(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	dashboard, err := h.dashboardService.UpdateDashboard(ctx, upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, dashboard); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeCreateCell(r *http.Request) (*manta.Cell, error) {
	var cell manta.Cell

	err := json.NewDecoder(r.Body).Decode(&cell)
	if err != nil {
		return nil, err
	}

	return &cell, nil
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

	if err = h.EncodeResponse(ctx, w, http.StatusCreated, cell); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardsHandler) getCell(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	dashboardID, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cellID, err := cellIDFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cell, err := h.dashboardService.FindDashboardCellByID(ctx, dashboardID, cellID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, cell); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeCellUpdate(r *http.Request) (*manta.DashboardCellUpdate, error) {
	var upd = &manta.DashboardCellUpdate{}

	dashboardID, err := idFromPath(r)
	if err != nil {
		return nil, err
	}

	cellID, err := cellIDFromPath(r)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		return nil, err
	}

	upd.DashboardID = dashboardID
	upd.CellID = cellID

	return upd, nil
}

func (h *DashboardsHandler) updateCell(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	upd, err := decodeCellUpdate(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cell, err := h.dashboardService.UpdateDashboardCell(ctx, *upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, cell); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func cellIDFromPath(r *http.Request) (manta.ID, error) {
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

	dashboard, err := h.dashboardService.FindDashboardByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusCreated, dashboard); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *DashboardsHandler) deleteCell(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	did, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cid, err := cellIDFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.dashboardService.RemoveDashboardCell(ctx, did, cid); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}

func (h *DashboardsHandler) changes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logs, _, err := findOplogByResourceID(r, h.operationLogService)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, logs); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
