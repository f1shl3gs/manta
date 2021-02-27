package web

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	NotificationEndpointPrefix = "/api/v1/notification_endpoint"
	NotificationEndpointIDPath = "/api/v1/notification_endpoint/:id"
)

type NotificationEndpointHandler struct {
	*Router

	logger                      *zap.Logger
	notificationEndpointService manta.NotificationEndpointService
}

func NewNotificationEndpointHandler(
	router *Router,
	logger *zap.Logger,
	nef manta.NotificationEndpointService,
) {
	h := &NotificationEndpointHandler{
		Router:                      router,
		logger:                      logger.With(zap.String("handler", "notification_endpoint")),
		notificationEndpointService: nef,
	}

	h.HandlerFunc(http.MethodGet, NotificationEndpointPrefix, h.handleList)
	h.HandlerFunc(http.MethodGet, NotificationEndpointIDPath, h.handleGet)
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

	id, err := idFromURI(r, "id")
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
