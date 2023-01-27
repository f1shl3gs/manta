package http

import (
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/http/router"

	"go.uber.org/zap"
)

type buildInfoHandler struct {
	*router.Router
}

type buildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Branch  string `json:"branch"`
}

func NewBuildInfoHandler(backend *Backend, logger *zap.Logger) {
	h := &buildInfoHandler{
		Router: backend.router,
	}

	info := &buildInfo{
		Version: manta.Version,
		Commit:  manta.Commit,
		Branch:  manta.Branch,
	}

	h.HandlerFunc(http.MethodGet, "/build", func(w http.ResponseWriter, r *http.Request) {
		err := h.EncodeResponse(r.Context(), w, http.StatusOK, info)
		if err != nil {
			logEncodingError(logger, r, err)
		}
	})
}
