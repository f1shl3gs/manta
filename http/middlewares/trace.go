package middlewares

import (
	"net/http"

	"github.com/f1shl3gs/manta/pkg/tracing"
	ua "github.com/mileusna/useragent"
)

func Trace(next http.Handler) http.Handler {
	name := "manta"
	fn := func(w http.ResponseWriter, r *http.Request) {
		span, r := tracing.ExtractFromHTTPRequest(r, name)
		defer span.Finish()

		span.LogKV("user_agent", UserAgent(r))
		for k, v := range r.Header {
			if len(v) == 0 {
				continue
			}

			if k == "Authorization" || k == "User-Agent" {
				continue
			}

			// If header has multiple values, only the first value will be logged on the traces.
			span.LogKV(k, v[0])
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func UserAgent(r *http.Request) string {
	header := r.Header.Get("User-Agent")
	if header == "" {
		return "unknown"
	}

	return ua.Parse(header).Name
}
