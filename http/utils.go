package http

import (
	"context"
	"github.com/f1shl3gs/manta"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func ExtractParamFromContext(ctx context.Context, name string) string {
	params := httprouter.ParamsFromContext(ctx)

	return params.ByName(name)
}

func OrgIdFromQuery(r *http.Request) (manta.ID, error) {
	var text = r.URL.Query().Get("orgId")
	return parseId(text)
}

func IDFromPath(r *http.Request) (manta.ID, error) {
	text := ExtractParamFromContext(r.Context(), "id")

	return parseId(text)
}

func parseId(s string) (id manta.ID, err error) {
	err = id.DecodeFromString(s)

	return id, err
}
