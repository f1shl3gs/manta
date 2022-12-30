package router

import (
	"net/http"

	"github.com/mileusna/useragent"
)

func UserAgent(r *http.Request) string {
	header := r.Header.Get("User-Agent")
	if header == "" {
		return "unknown"
	}

	return useragent.Parse(header).Name
}
