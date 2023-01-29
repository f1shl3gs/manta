package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"
	"github.com/f1shl3gs/manta/notification"
)

const (
	notificationEndpointPrefix = apiV1Prefix + "/notificationEndpoints"
	notificationEndpointIDPath = notificationEndpointPrefix + "/:id"
)

type NotificationEndpointHandler struct {
	*router.Router
	logger *zap.Logger

	notificationEndpointService manta.NotificationEndpointService
}

func NewNotificationEendpointHandler(logger *zap.Logger, backend *Backend) {
	h := NotificationEndpointHandler{
		Router:                      backend.router,
		logger:                      logger.With(zap.String("handler", "notification_endpoint")),
		notificationEndpointService: backend.NotificationEndpointService,
	}

	h.HandlerFunc(http.MethodGet, notificationEndpointPrefix, h.handleList)
	h.HandlerFunc(http.MethodPost, notificationEndpointPrefix, h.handleCreate)
	h.HandlerFunc(http.MethodGet, notificationEndpointIDPath, h.handleGet)
	h.HandlerFunc(http.MethodPatch, notificationEndpointIDPath, h.handlePatch)
	h.HandlerFunc(http.MethodPost, notificationEndpointIDPath, h.handleUpdate)
	h.HandlerFunc(http.MethodDelete, notificationEndpointIDPath, h.handleDelete)
}

func (h *NotificationEndpointHandler) handleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgID, err := orgIDFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	list, err := h.notificationEndpointService.FindNotificationEndpoints(ctx, manta.NotificationEndpointFilter{
		OrgID: orgID,
	})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, list); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeNotificationEndpoint(r *http.Request) (manta.NotificationEndpoint, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return notification.UnmarshalJSON(data)
}

func (h *NotificationEndpointHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	ne, err := decodeNotificationEndpoint(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.notificationEndpointService.CreateNotificationEndpoint(ctx, ne)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusCreated, ne); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *NotificationEndpointHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	ne, err := h.notificationEndpointService.FindNotificationEndpointByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, ne); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeNotificationEndpointUpdate(r *http.Request) (manta.NotificationEndpointUpdate, error) {
	var upd manta.NotificationEndpointUpdate

	err := json.NewDecoder(r.Body).Decode(&upd)
	if err != nil {
		return manta.NotificationEndpointUpdate{}, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "invalid notification endpoint update",
			Err:  err,
		}
	}

	return upd, nil
}

func (h *NotificationEndpointHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	upd, err := decodeNotificationEndpointUpdate(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	ne, err := h.notificationEndpointService.PatchNotificationEndpoint(ctx, id, upd)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, ne); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *NotificationEndpointHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	ne, err := decodeNotificationEndpoint(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	ne, err = h.notificationEndpointService.UpdateNotificationEndpoint(ctx, id, ne)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, ne); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *NotificationEndpointHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	_, _, err = h.notificationEndpointService.DeleteNotificationEndpoint(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}
