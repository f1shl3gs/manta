package http

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
)

const (
	configurationPrefix = apiV1Prefix + `/configurations`
	configurationWithID = configurationPrefix + `/:id`
)

type ConfigurationHandler struct {
	*Router

	logger               *zap.Logger
	configurationService manta.ConfigurationService
}

func NewConfigurationService(backend *Backend, logger *zap.Logger) {
	h := &ConfigurationHandler{
		Router:               backend.router,
		logger:               logger.With(zap.String("handle", "configuration")),
		configurationService: backend.ConfigurationService,
	}

	h.HandlerFunc(http.MethodGet, configurationPrefix, h.listConfigurations)
	h.HandlerFunc(http.MethodPost, configurationPrefix, h.createConfiguration)
	h.HandlerFunc(http.MethodGet, configurationWithID, h.getConfiguration)
	h.HandlerFunc(http.MethodPatch, configurationWithID, h.updateConfiguration)
	h.HandlerFunc(http.MethodDelete, configurationWithID, h.deleteConfiguration)
}

func (h *ConfigurationHandler) listConfigurations(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	orgID, err := orgIdFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cs, err := h.configurationService.FindConfigurations(ctx, manta.ConfigurationFilter{OrgID: orgID})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, cs); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ConfigurationHandler) getConfiguration(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	c, err := h.configurationService.GetConfiguration(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusOK, c); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ConfigurationHandler) updateConfiguration(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		upd = manta.ConfigurationUpdate{}
	)

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&upd); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.configurationService.UpdateConfiguration(ctx, id, upd); err != nil {
		h.HandleHTTPError(ctx, err, w)
	}
}

func (h *ConfigurationHandler) deleteConfiguration(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.configurationService.DeleteConfiguration(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}

func (h *ConfigurationHandler) createConfiguration(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		c   = manta.Configuration{}
	)

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.configurationService.CreateConfiguration(ctx, &c); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = encodeResponse(ctx, w, http.StatusCreated, &c); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
