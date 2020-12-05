package web

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

const (
	LokiPrefix = "/loki/api/v1"
)

type lokiHandler struct {
	logger *zap.Logger
	proxy  *httputil.ReverseProxy
}

func (h *lokiHandler) handle(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}

func newLokiHandler(logger *zap.Logger, backend string) (*lokiHandler, error) {
	target, err := url.Parse(backend)
	if err != nil {
		return nil, err
	}

	h := &lokiHandler{
		logger: logger,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}

	return h, nil
}

func NewLokiService(h *lokiHandler, router *Router) {
	router.HandlerFunc(http.MethodGet, LokiPrefix+"/*dummy", h.handle)
}
