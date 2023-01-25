package http

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"
)

const (
	checksPrefix     = apiV1Prefix + `/checks`
	checksIDPath     = checksPrefix + `/:id`
	checksChangePath = checksIDPath + `/changes`
)

type ChecksHandler struct {
	*router.Router

	logger       *zap.Logger
	checkService manta.CheckService
	taskService  manta.TaskService
	oplogService manta.OperationLogService
}

func NewChecksHandler(logger *zap.Logger, router *router.Router, cs manta.CheckService, ts manta.TaskService, ol manta.OperationLogService) {
	h := &ChecksHandler{
		Router:       router,
		logger:       logger.With(zap.String("handler", "check")),
		checkService: cs,
		taskService:  ts,
		oplogService: ol,
	}

	h.HandlerFunc(http.MethodGet, checksPrefix, h.handleList)
	h.HandlerFunc(http.MethodPost, checksPrefix, h.handleCreate)
	h.HandlerFunc(http.MethodDelete, checksIDPath, h.handleDelete)
	h.HandlerFunc(http.MethodPost, checksIDPath, h.handleUpdate)
	h.HandlerFunc(http.MethodPatch, checksIDPath, h.handlePatch)
	h.HandlerFunc(http.MethodGet, checksIDPath, h.handleGet)
	h.HandlerFunc(http.MethodGet, checksChangePath, h.handleChanges)
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

	orgID, err := orgIdFromQuery(r)
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

	err = h.EncodeResponse(ctx, w, http.StatusOK, &list)
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
			Msg:  "decode check failed",
			Err:  err,
		}
	}

	err = c.Validate()
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "validate check failed",
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

	if err = h.EncodeResponse(ctx, w, http.StatusCreated, c); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ChecksHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromPath(r)
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
	upd := manta.CheckUpdate{}

	err := json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		return upd, &manta.Error{Code: manta.EInvalid, Op: "decode CheckUpdate", Err: err}
	}

	return upd, nil
}

func (h *ChecksHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	check, err := decodeCheck(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	check, err = h.checkService.UpdateCheck(ctx, id, check)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err := h.EncodeResponse(ctx, w, http.StatusOK, check); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ChecksHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	c, err := h.checkService.FindCheckByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, c); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ChecksHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	upd, err := decodeCheckUpdate(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	check, err := h.checkService.PatchCheck(ctx, id, upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
	if err = h.EncodeResponse(ctx, w, http.StatusOK, check); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ChecksHandler) handleChanges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logs, _, err := findOplogByResourceID(r, h.oplogService)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, logs); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
