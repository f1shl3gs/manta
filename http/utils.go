package http

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/f1shl3gs/manta"
)

func extractParamFromContext(ctx context.Context, name string) string {
	params := httprouter.ParamsFromContext(ctx)

	return params.ByName(name)
}

func orgIdFromQuery(r *http.Request) (manta.ID, error) {
	var (
		text = r.URL.Query().Get("orgId")
		id   manta.ID
	)

	return id, id.DecodeFromString(text)
}

func idFromPath(r *http.Request) (manta.ID, error) {
	var (
		text = extractParamFromContext(r.Context(), "id")
		id   manta.ID
	)

	return id, id.DecodeFromString(text)
}
