package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*httprouter.Router
	manta.HTTPErrorHandler

	noAuthRouter *httprouter.Router
}

func NewRouter() *Router {
	router := httprouter.New()
	router.SaveMatchedRoutePath = true
	router.NotFound = http.NotFoundHandler()

	return &Router{
		Router: router,
	}
}

// HandleHTTPError implement HTTPErrorHandler
func (r *Router) HandleHTTPError(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		return
	}

	code := manta.ErrorCode(err)
	w.Header().Set(PlatformErrorCodeHeader, code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(ErrorCodeToStatusCode(ctx, code))
	var e struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	e.Code = manta.ErrorCode(err)
	if mErr, ok := err.(*manta.Error); ok {
		e.Message = mErr.Error()
	} else {
		e.Message = "An internal error has occurred, " + err.Error()
	}
	b, _ := json.Marshal(e)
	_, _ = w.Write(b)
}

func (r *Router) RegisterNoAuthRoute(method, path string) {
	r.noAuthRouter.HandlerFunc(method, path, func(w http.ResponseWriter, r *http.Request) {})
}
