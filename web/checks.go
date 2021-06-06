package web

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	ChecksPrefix = `/api/v1/checks`
	ChecksIDPath = `/api/v1/checks/:id`
)

type ChecksHandler struct {
	*Router

	logger       *zap.Logger
	checkService manta.CheckService
	taskService  manta.TaskService
}

func NewChecksHandler(logger *zap.Logger, router *Router, cs manta.CheckService, ts manta.TaskService) {
	h := &ChecksHandler{
		Router:       router,
		logger:       logger.With(zap.String("handler", "check")),
		checkService: cs,
		taskService:  ts,
	}

	h.HandlerFunc(http.MethodGet, ChecksPrefix, h.handleList)
	h.HandlerFunc(http.MethodPut, ChecksPrefix, h.handleCreate)
	h.HandlerFunc(http.MethodDelete, ChecksIDPath, h.handleDelete)
	h.HandlerFunc(http.MethodPost, ChecksIDPath, h.handleUpdate)
	h.HandlerFunc(http.MethodPatch, ChecksIDPath, h.handlePatch)
	h.HandlerFunc(http.MethodGet, ChecksIDPath, h.handleGet)
}

type check struct {
	*manta.Check

	LatestCompleted time.Time `json:"latestCompleted"`
	LatestScheduled time.Time `json:"latestScheduled"`
	LatestSuccess   time.Time `json:"latestSuccess"`
	LatestFailure   time.Time `json:"latestFailure"`
	LastRunStatus   string    `json:"lastRunStatus"`
	LastRunError    string    `json:"lastRunError,omitempty"`
}

func (h *ChecksHandler) handleList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		filter = manta.CheckFilter{}
	)

	orgID, err := orgIDFromRequest(r)
	if err == nil {
		filter.OrgID = &orgID
	}

	checks, _, err := h.checkService.FindChecks(ctx, filter)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	list := make([]check, 0, len(checks))
	for _, c := range checks {
		task, err := h.taskService.FindTaskByID(ctx, c.TaskID)
		if err != nil {
			h.logger.Warn("Find task by id failed",
				zap.String("task", c.TaskID.String()),
				zap.Error(err))

			continue
		}

		list = append(list, check{
			Check:           c,
			LatestCompleted: task.LatestCompleted,
			LatestScheduled: task.LatestScheduled,
			LatestSuccess:   task.LatestSuccess,
			LatestFailure:   task.LatestFailure,
			LastRunStatus:   task.LastRunStatus,
			LastRunError:    task.LastRunError,
		})
	}

	err = encodeResponse(ctx, w, http.StatusOK, &list)
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeCheck(r *http.Request) (*manta.Check, error) {
	c := &manta.Check{}
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Op:   "Decode check",
			Err:  err,
		}
	}

	err = c.Validate()
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Op:   "Validate check",
			Err:  err,
		}
	}

	return c, nil
}

func (h *ChecksHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	c, err := decodeCheck(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.checkService.CreateCheck(ctx, c)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ChecksHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.checkService.DeleteCheck(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeCheckUpdate(r *http.Request) (manta.CheckUpdate, error) {
	udp := manta.CheckUpdate{}

	err := json.NewDecoder(r.Body).Decode(&udp)
	if err != nil {
		return udp, &manta.Error{Code: manta.EInvalid, Op: "decode CheckUpdate", Err: err}
	}

	if err = udp.Validate(); err != nil {
		return udp, err
	}

	return udp, nil
}

func (h *ChecksHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	c, err := decodeCheck(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	_, err = h.checkService.UpdateCheck(ctx, id, c)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChecksHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	c, err := h.checkService.FindCheckByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, c); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ChecksHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	udp, err := decodeCheckUpdate(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	_, err = h.checkService.PatchCheck(ctx, id, udp)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
