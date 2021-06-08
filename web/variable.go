package web

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	VariablePrefix = "/api/v1/orgs/:orgID/variables"
	VariableIDPath = "/api/v1/orgs/:orgID/variables/:id"
)

type VariableHandler struct {
	*Router

	logger          *zap.Logger
	variableService manta.VariableService
}

func NewVariableHandler(logger *zap.Logger, router *Router, variableService manta.VariableService) {
	h := &VariableHandler{
		Router:          router,
		logger:          logger.With(zap.String("handler", "variable")),
		variableService: variableService,
	}

	h.HandlerFunc(http.MethodGet, VariablePrefix, h.handleList)
	h.HandlerFunc(http.MethodPost, VariablePrefix, h.handleCreate)
	h.HandlerFunc(http.MethodGet, VariableIDPath, h.handleGet)
	h.HandlerFunc(http.MethodDelete, VariableIDPath, h.handleDelete)
	h.HandlerFunc(http.MethodPatch, VariableIDPath, h.handlePatch)
	h.HandlerFunc(http.MethodPost, VariableIDPath, h.handleUpdate)
}

// handleList is the http handler for List variables
func (h *VariableHandler) handleList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	vars, err := h.variableService.FindVariables(ctx, manta.VariableFilter{OrgID: &orgID})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, vars); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

// handleGet returns the variable find by id
func (h *VariableHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	variable, err := h.variableService.FindVariableByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, variable); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

// handleDelete deletes variable by id
func (h *VariableHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.variableService.DeleteVariable(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeVariableUpdate(r *http.Request) (*manta.VariableUpdate, error) {
	upd := &manta.VariableUpdate{}

	err := json.NewDecoder(r.Body).Decode(upd)
	if err != nil {
		return nil, err
	}

	return upd, nil
}

// handlePatch update a variable with changeset
func (h *VariableHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	upd, err := decodeVariableUpdate(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	_, err = h.variableService.PatchVariable(ctx, id, upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func decodeVariable(r *http.Request) (*manta.Variable, error) {
	variable := &manta.Variable{}

	err := json.NewDecoder(r.Body).Decode(variable)
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "Decode variable failed",
			Err:  err,
		}
	}

	err = variable.Validate()
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "Validate variable failed",
			Err:  err,
		}
	}

	return variable, nil
}

func (h *VariableHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	variable, err := decodeVariable(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.variableService.CreateVariable(ctx, variable)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *VariableHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	variable, err := decodeVariable(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.variableService.UpdateVariable(ctx, variable)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
