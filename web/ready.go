package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func ReadyHandler() http.Handler {
	up := time.Now()
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		var status = struct {
			Status string    `json:"status"`
			Start  time.Time `json:"started"`
			// TODO(jsteenb2): learn why and leave comment for this being a toml.Duration
			Up time.Duration `json:"up"`
		}{
			Status: "ready",
			Start:  up,
			Up:     time.Since(up),
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		if err := enc.Encode(status); err != nil {
			fmt.Fprintf(w, "Error encoding status data: %v\n", err)
		}
	}

	return http.HandlerFunc(fn)
}
