package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"
	"github.com/f1shl3gs/manta/raftstore"
	raftstorepb "github.com/f1shl3gs/manta/raftstore/pb"

	"go.uber.org/zap"
)

const (
	clusterServicePrefix = apiV1Prefix + "/cluster"
	clusterServiceIDPath = clusterServicePrefix + "/:id"
)

type ClusterServiceHandler struct {
	*router.Router

	logger      *zap.Logger
	raftService raftstore.ClusterService
}

func NewClusterServiceHandler(logger *zap.Logger, backend *Backend) {
	if backend.ClusterService == raftstore.ClusterService(nil) {
		return
	}

	h := &ClusterServiceHandler{
		Router:      backend.router,
		logger:      logger,
		raftService: backend.ClusterService,
	}

	h.HandlerFunc(http.MethodGet, clusterServicePrefix, h.list)
	h.HandlerFunc(http.MethodPost, clusterServicePrefix, h.add)
	h.HandlerFunc(http.MethodDelete, clusterServiceIDPath, h.delete)
}

func (h *ClusterServiceHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	members := h.raftService.Members()

	if err := h.EncodeResponse(ctx, w, http.StatusOK, members); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ClusterServiceHandler) add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	member := raftstorepb.Member{}

	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		h.HandleHTTPError(ctx, &manta.Error{
			Code: manta.EInvalid,
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

func (h *ClusterServiceHandler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	text := extractParamFromContext(r.Context(), "id")
	id, err := strconv.ParseUint(text, 16, 64)
	if err != nil {
		h.HandleHTTPError(ctx, &manta.Error{Code: manta.EInvalid, Msg: "invalid node id", Err: err}, w)
		return
	}

	err = h.raftService.Remove(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}
