package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/f1shl3gs/manta"
)

func extractParamFromContext(ctx context.Context, name string) string {
	params := httprouter.ParamsFromContext(ctx)

	return params.ByName(name)
}

func orgIDFromQuery(r *http.Request) (manta.ID, error) {
	var (
		text = r.URL.Query().Get("orgID")
		id   manta.ID
	)

	err := id.DecodeFromString(text)
	if err != nil {
		return 0, &manta.Error{
			Code: manta.EInvalid,
			Msg:  "invalid organization id found in query",
			Err:  err,
		}
	}

	return id, nil
}

func idFromPath(r *http.Request) (manta.ID, error) {
	var (
		text = extractParamFromContext(r.Context(), "id")
		id   manta.ID
	)

	return id, id.DecodeFromString(text)
}

func limitFromQuery(r *http.Request, defaultValue, max int64) (int, error) {
	var (
		n   int64
		err error
	)

	if max == 0 {
		max = 500
	}

	if defaultValue == 0 {
		defaultValue = max
	}

	s := r.URL.Query().Get("limit")
	if s == "" {
		n = defaultValue
	} else {
		n, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, &manta.Error{
				Code: manta.EInvalid,
				Msg:  "Parse limit failed",
				Op:   "parse limit",
				Err:  err,
			}
		}
	}

	if n > max || n <= 0 {
		return 0, &manta.Error{
			Code: manta.EUnprocessableEntity,
			Msg:  fmt.Sprintf("Limit value must between 1 and %d", max),
		}
	}

	return int(n), err
}
