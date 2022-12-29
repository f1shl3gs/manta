package router

import (
	"net/http"

	"github.com/f1shl3gs/manta/pkg/tracing"
)

func Trace() Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            span, r := tracing.ExtractFromHTTPRequest(r, "manta")
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
    }
}