package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/config"
	"github.com/f1shl3gs/manta/http/router"
)

const (
	configPrefix = apiV1Prefix + `/configs`
	configWithID = configPrefix + `/:id`
)

var (
	watchStreams = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "config",
		Name:      "watching_streams_total",
	})
)

type ConfigHandler struct {
	*router.Router

	logger        *zap.Logger
	configService *config.CoordinatingConfigService
}

func NewConfigService(backend *Backend, logger *zap.Logger) {
	h := &ConfigHandler{
		Router:        backend.router,
		logger:        logger.With(zap.String("handle", "config")),
		configService: config.NewCoordinatingVertexService(backend.ConfigService, logger),
	}

	backend.PromRegistry.MustRegister(watchStreams)

	h.HandlerFunc(http.MethodGet, configPrefix, h.listConfigs)
	h.HandlerFunc(http.MethodPost, configPrefix, h.createConfig)
	h.HandlerFunc(http.MethodGet, configWithID, h.getConfig)
	h.HandlerFunc(http.MethodPatch, configWithID, h.updateConfig)
	h.HandlerFunc(http.MethodDelete, configWithID, h.deleteConfig)
}

func (h *ConfigHandler) listConfigs(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	orgID, err := orgIDFromQuery(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	cs, err := h.configService.FindConfigs(ctx, manta.ConfigFilter{OrgID: orgID})
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusOK, cs); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *ConfigHandler) getConfig(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	var encodeResp = func(cf *manta.Config, writer io.Writer) error {
		var (
			data []byte
			err  error
		)

		if strings.Contains(r.Header.Get("accept"), "json") {
			data, err = json.Marshal(cf)
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
		cf, err := h.configService.FindConfigByID(ctx, id)
		if err != nil {
			h.HandleHTTPError(ctx, err, w)
			return
		}

		if err = encodeResp(cf, w); err != nil {
			logEncodingError(h.logger, r, err)
		}

		return
	}

	watchStreams.Inc()
	defer watchStreams.Dec()

	// watching first then get, so we won't miss any updates (no promise)
	queue := h.configService.Sub(id)
	defer queue.Close()

	first, err := h.configService.FindConfigByID(ctx, id)
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
		cf     *manta.Config
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

func (h *ConfigHandler) updateConfig(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		upd = manta.ConfigUpdate{}
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

	if cf, err := h.configService.UpdateConfig(ctx, id, upd); err != nil {
		h.HandleHTTPError(ctx, err, w)
	} else {
		if err = h.EncodeResponse(ctx, w, http.StatusOK, cf); err != nil {
			logEncodingError(h.logger, r, err)
		}
	}
}

func (h *ConfigHandler) deleteConfig(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	id, err := idFromPath(r)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = h.configService.DeleteConfig(ctx, id)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}
}

func (h *ConfigHandler) createConfig(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		c   = manta.Config{}
	)

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.configService.CreateConfig(ctx, &c); err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	if err = h.EncodeResponse(ctx, w, http.StatusCreated, &c); err != nil {
		logEncodingError(h.logger, r, err)
	}
}
