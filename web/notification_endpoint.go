package web

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	NotificationEndpointPrefix = "/api/v1/orgs/:orgID/notification_endpoints"
	NotificationEndpointIDPath = "/api/v1/orgs/:orgID/notification_endpoints/:id"
)

type NotificationEndpointHandler struct {
	*Router

	logger                      *zap.Logger
	notificationEndpointService manta.NotificationEndpointService
}

func NewNotificationEndpointHandler(
	logger *zap.Logger,
	router *Router,
	nef manta.NotificationEndpointService,
) {
	h := &NotificationEndpointHandler{
		Router:                      router,
		logger:                      logger.With(zap.String("handler", "notification_endpoint")),
		notificationEndpointService: nef,
	}

	h.HandlerFunc(http.MethodPost, NotificationEndpointPrefix, h.handleCreate)
	h.HandlerFunc(http.MethodGet, NotificationEndpointPrefix, h.handleList)
	h.HandlerFunc(http.MethodGet, NotificationEndpointIDPath, h.handleGet)
	h.HandlerFunc(http.MethodPatch, NotificationEndpointIDPath, h.handlePatch)
	h.HandlerFunc(http.MethodDelete, NotificationEndpointIDPath, h.handleDelete)
}

func decodeNotificationEndpoint(r *http.Request) (*manta.NotificationEndpoint, error) {
	endpoint := &manta.NotificationEndpoint{}

	err := json.NewDecoder(r.Body).Decode(endpoint)
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "decode notification endpoint failed",
			Err:  err,
		}
	}

	err = endpoint.Validate()
	if err != nil {
		return nil, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "validate failed",
			Err:  err,
		}
	}

	return endpoint, nil
}

// handleCreate is the http handler for create a new notification endpoint
func (h *NotificationEndpointHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	endpoint, err := decodeNotificationEndpoint(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.notificationEndpointService.CreateNotificationEndpoint(ctx, endpoint)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// handleList is the http handler for list all notification endpoints
func (h *NotificationEndpointHandler) handleList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	orgID, err := orgIDFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	nes, _, err := h.notificationEndpointService.FindNotificationEndpoints(ctx, manta.NotificationEndpointFilter{OrgID: &orgID})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, &nes); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *NotificationEndpointHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	ne, err := h.notificationEndpointService.FindNotificationEndpointByID(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, ne); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *NotificationEndpointHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, err := idFromRequest(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.notificationEndpointService.DeleteNotificationEndpoint(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func decodeNotificationEndpointUpdate(r *http.Request) (manta.ID, manta.NotificationEndpointUpdate, error) {
	id, err := idFromRequest(r)
	if err != nil {
		return 0, manta.NotificationEndpointUpdate{}, err
	}

	udp := manta.NotificationEndpointUpdate{}
	err = json.NewDecoder(r.Body).Decode(&udp)
	if err != nil {
		return 0, manta.NotificationEndpointUpdate{}, err
	}

	return id, udp, nil
}

func (h *NotificationEndpointHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	id, udp, err := decodeNotificationEndpointUpdate(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	_, err = h.notificationEndpointService.UpdateNotificationEndpoint(ctx, id, udp)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
