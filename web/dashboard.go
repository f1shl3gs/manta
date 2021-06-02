package web

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	DashboardPrefix     = "/api/v1/dashboards"
	DashboardIDPath     = "/api/v1/dashboards/:id"
	DashboardCellPrefix = "/api/v1/dashboards/:id/cells"
	DashboardCellIDPath = "/api/v1/dashboards/:id/cells/:cellID"
)

type DashboardHandler struct {
	*Router

	logger           *zap.Logger
	dashboardService manta.DashboardService
}

func NewDashboardService(h *DashboardHandler) {
	h.HandlerFunc(http.MethodGet, DashboardPrefix, h.handleList)
	h.HandlerFunc(http.MethodGet, DashboardIDPath, h.handleGet)
	h.HandlerFunc(http.MethodDelete, DashboardIDPath, h.handleDelete)
	h.HandlerFunc(http.MethodPost, DashboardPrefix, h.handleCreate)
	h.HandlerFunc(http.MethodPost, DashboardCellPrefix, h.handleAddCell)
	h.HandlerFunc(http.MethodGet, DashboardCellIDPath, h.handleGetCell)
	h.HandlerFunc(http.MethodPut, DashboardCellPrefix, h.handleReplaceDashboardCells)
	h.HandlerFunc(http.MethodPatch, DashboardIDPath, h.handleUpdate)
	h.HandlerFunc(http.MethodPatch, DashboardCellIDPath, h.handleUpdateCell)
	h.HandlerFunc(http.MethodDelete, DashboardCellIDPath, h.handleDeleteCell)
}

func (h *DashboardHandler) handleList(w http.ResponseWriter, r *http.Request) {
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

func (h *DashboardHandler) handleGet(w http.ResponseWriter, r *http.Request) {
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

func (h *DashboardHandler) handleReplaceDashboardCells(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromRequestPath(r)
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

func (h *DashboardHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusNoContent)
}

func (h *DashboardHandler) handleGetCell(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		params = httprouter.ParamsFromContext(ctx)
	)

	dashID, err := idFromParams(params, "id")
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cellID, err := idFromParams(params, "cellID")
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cell, err := h.dashboardService.GetDashboardCell(ctx, dashID, cellID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, cell); err != nil {
		h.HandleHTTPError(ctx, err, w)
	}
}

func decodeDashboardUpdate(r *http.Request) (manta.DashboardUpdate, error) {
	var udp manta.DashboardUpdate

	err := json.NewDecoder(r.Body).Decode(&udp)
	if err != nil {
		return manta.DashboardUpdate{}, errors.Wrap(err, "decode DashboardUpdate failed")
	}

	return udp, nil
}

func (h *DashboardHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := idFromRequestPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	udp, err := decodeDashboardUpdate(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	dash, err := h.dashboardService.UpdateDashboard(ctx, id, udp)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, dash); err != nil {
		logEncodingError(h.logger, r, err)
		return
	}
}

func (h *DashboardHandler) handleUpdateCell(w http.ResponseWriter, r *http.Request) {
	var (
		udp manta.DashboardCellUpdate
		ctx = r.Context()
	)

	dashboardID, err := idFromURI(r, "id")
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cellID, err := idFromURI(r, "cellID")
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&udp)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	_, err = h.dashboardService.UpdateDashboardCell(ctx, dashboardID, cellID, udp)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DashboardHandler) handleDeleteCell(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	dashboardID, err := idFromURI(r, "id")
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cellID, err := idFromURI(r, "cellID")
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.dashboardService.RemoveDashboardCell(ctx, dashboardID, cellID)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
