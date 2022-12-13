package http

import (
	"encoding/json"
    "io"
    "net/http"
	"net/http/httputil"
    "strings"

    "go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/vertex"
)

const (
	configurationPrefix = apiV1Prefix + `/configurations`
	configurationWithID = configurationPrefix + `/:id`
)

type ConfigurationHandler struct {
	*Router

	logger               *zap.Logger
	configurationService *vertex.CoordinatingVertexService
}

func NewConfigurationService(backend *Backend, logger *zap.Logger) {
	h := &ConfigurationHandler{
		Router:               backend.router,
		logger:               logger.With(zap.String("handle", "configuration")),
		configurationService: vertex.NewCoordinatingVertexService(backend.ConfigurationService, logger),
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

    var encodeResp = func(cf *manta.Configuration, writer io.Writer) error {
        var (
            data []byte
            err error
        )

        if strings.Contains(r.Header.Get("accept"), "json") {
            data, err = cf.Marshal()
            if err != nil {
                return err
            }
        } else {
            data = []byte(cf.Data)
        }

        _, err = writer.Write(data)
        return err
    }

	if r.URL.Query().Get("watch") != "true" {
		// Just get config, not watching
		cf, err := h.configurationService.GetConfiguration(ctx, id)
		if err != nil {
			h.HandleHTTPError(ctx, err, w)
			return
		}

        if err = encodeResp(cf, w); err != nil {
			logEncodingError(h.logger, r, err)
		}

		return
	}

	// watching first then get, so we won't miss any updates (no promise)
	queue := h.configurationService.Sub(id)
	defer queue.Close()

	first, err := h.configurationService.GetConfiguration(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	var (
		cf     *manta.Configuration
		closed bool
		writer = httputil.NewChunkedWriter(w)
	)

	defer writer.Close()

	for {
		if first != nil {
			cf = first
			first = nil
		}

        // what we watched for is configuratin, not the data field,
        // so false notification might happenned.
        err = encodeResp(cf, writer)
		if err != nil {
			h.logger.Warn("watch failed",
				zap.String("client", r.RemoteAddr),
				zap.Error(err))
			return
		}

		flusher.Flush()

		// wait for new channel
		select {
		case <-ctx.Done():
			return

		case cf, closed = <-queue.Ch():
			if !closed {
				// config deleted, so the channel is closed too
				return
			}
		}
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

	if cf, err := h.configurationService.UpdateConfiguration(ctx, id, upd); err != nil {
		h.HandleHTTPError(ctx, err, w)
	} else {
		if err = encodeResponse(ctx, w, http.StatusOK, cf); err != nil {
			logEncodingError(h.logger, r, err)
		}
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
