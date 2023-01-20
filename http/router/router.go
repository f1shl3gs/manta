package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/f1shl3gs/manta"
)

// PlatformErrorCodeHeader shows the error code of platform error.
const PlatformErrorCodeHeader = "X-Platform-Error-Code"

type Middleware func(next http.HandlerFunc) http.HandlerFunc

type Router struct {
	router      *httprouter.Router
	middlewares []Middleware
}

func New(middlewares ...Middleware) *Router {
	router := &httprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		SaveMatchedRoutePath:   true,
	}

	return &Router{
		router:      router,
		middlewares: middlewares,
	}
}

// ServeHTTP implement http.Handle
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) Use(ms ...Middleware) {
	r.middlewares = append(r.middlewares, ms...)
}

func (r *Router) HandlerFunc(method, path string, handlerFunc http.HandlerFunc) {
	handler := handlerFunc
	for _, m := range r.middlewares {
		handler = m(handler)
	}

	r.router.HandlerFunc(method, path, handler)
}

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

func (r *Router) EncodeResponse(ctx context.Context, w http.ResponseWriter, status int, payload interface{}) error {
	if payload == nil {
		w.WriteHeader(status)
		return nil
	}

	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return err
	}

	return nil
}
