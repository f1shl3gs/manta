package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/f1shl3gs/manta/errors"
	"github.com/f1shl3gs/manta/http/router"
	"github.com/f1shl3gs/manta/raftstore"

	"go.uber.org/zap"
)

const (
	raftServicePrefix = apiV1Prefix + "/cluster"
	raftServiceIDPath = raftServicePrefix + "/:id"
)

type RaftServiceHandler struct {
	*router.Router

	logger      *zap.Logger
	raftService raftstore.RaftService
}

func NewRaftServiceHandler(logger *zap.Logger, backend *Backend) {
	if backend.RaftService == nil {
		return
	}

	h := &RaftServiceHandler{
		Router:      backend.router,
		logger:      logger,
		raftService: backend.RaftService,
	}

	h.HandlerFunc(http.MethodPost, raftServicePrefix, h.add)
	h.HandlerFunc(http.MethodDelete, raftServiceIDPath, h.delete)
}

func (h *RaftServiceHandler) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	member := raftstore.Member{}

	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		h.HandleHTTPError(ctx, &errors.Error{
			Code: errors.EInvalid,
			Msg:  "decode cluster member failed",
			Err:  err,
		}, w)
		return
	}

	err = h.raftService.Add(ctx, member)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *RaftServiceHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	text := extractParamFromContext(r.Context(), "id")
	id, err := strconv.ParseUint(text, 16, 64)
	if err != nil {
		h.HandleHTTPError(ctx, &errors.Error{Code: errors.EInvalid, Msg: "invalid node id", Err: err}, w)
		return
	}

	err = h.raftService.Remove(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}
