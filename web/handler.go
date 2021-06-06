package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/julienschmidt/httprouter"

	"go.uber.org/zap"
)

var (
	ErrParamsNotFound = errors.New("params not found")
)

func idFromRequest(r *http.Request) (manta.ID, error) {
	var (
		id     manta.ID
		params = httprouter.ParamsFromContext(r.Context())
	)

	if err := id.DecodeFromString(params.ByName("id")); err != nil {
		return 0, err
	}

	return id, nil
}

func orgIDFromRequest(r *http.Request) (manta.ID, error) {
	var (
		id  manta.ID
		txt = r.URL.Query().Get("orgID")
	)

	err := id.DecodeFromString(txt)
	if err != nil {
		return 0, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "invalid orgID from url",
			Op:   "decode orgID",
			Err:  err,
		}
	}

	return id, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(res)
}

func logEncodingError(log *zap.Logger, r *http.Request, err error) {
	// If we encounter an error while encoding the response to an http request
	// the best thing we can do is log that error, as we may have already written
	// the headers for the http request in question.
	log.Info("Error encoding response",
		zap.String("path", r.URL.Path),
		zap.String("method", r.Method),
		zap.String("remote", r.RemoteAddr),
		zap.Error(err))
}
