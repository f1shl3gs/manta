package web

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// TODO: For now, this implement is just a demo, and all data is store in memory.

const (
	VertexPrefix = "/api/v1/vertex"
)

type instance struct {
	attributes    map[string]string
	lastHeartbeat time.Time
}

type VertexHandler struct {
	*Router

	logger *zap.Logger

	mtx       sync.RWMutex
	instances map[string]*instance
}

func NewVertexHandler(logger *zap.Logger, router *Router) {
	h := &VertexHandler{
		Router:    router,
		logger:    logger,
		instances: make(map[string]instance),
	}

	h.HandlerFunc(http.MethodGet, VertexPrefix, h.handleGet)
}

var upgrader = websocket.Upgrader{}

func (h *VertexHandler) listInstances(w http.ResponseWriter, r *http.Request) {
	h.mtx.RLock()
	list := make([]*instance, 0, len(h.instances))
	for _, ins := range h.instances {
		list = append(list, ins)
	}
	h.mtx.RUnlock()

	if err := encodeResponse(r.Context(), w, http.StatusOK, list); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *VertexHandler) handleGet(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Warn("Websocket upgrade failed", zap.Error(err))
		return
	}

	defer conn.Close()

	msg := struct {
		Type       string            `json:"type"`
		UUID       string            `json:"uuid"`
		Attributes map[string]string `json:"attributes"`
	}{}

	// TODO: start a goroutine for watch configuration and send to vertex

	for {
		if err = conn.ReadJSON(&msg); err != nil {
			h.logger.Warn("Read message from websocket failed",
				zap.Error(err))
			return
		}

		switch msg.Type {
		case "heartbeat":
			h.mtx.Lock()
			ins, exist := h.instances[msg.UUID]
			if !exist {
				ins = &instance{
					attributes: msg.Attributes,
				}

				h.instances[msg.UUID] = ins
			}

			ins.lastHeartbeat = time.Now()

			h.mtx.Unlock()
		default:
			h.logger.Warn("unknown message type from vertex",
				zap.String("type", msg.Type))
			return
		}
	}
}
