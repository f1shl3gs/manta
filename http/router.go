package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// PlatformErrorCodeHeader shows the error code of platform error.
const PlatformErrorCodeHeader = "X-Platform-Error-Code"

type Router struct {
	*httprouter.Router

	logger *zap.Logger
}

func (h *Router) HandleHTTPError(ctx context.Context, err error, w http.ResponseWriter) {
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

func encodeResponse(ctx context.Context, w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if payload != nil {
		return json.NewEncoder(w).Encode(payload)
	}

	return nil
}
