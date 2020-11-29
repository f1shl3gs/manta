package authorizer

import (
	"context"

	"github.com/f1shl3gs/manta"
)

type contextKey string

const (
	authorizerCtxKey contextKey = "authorizer"
)

func SetAuthorizer(ctx context.Context, a manta.Authorizer) context.Context {
	return context.WithValue(ctx, authorizerCtxKey, a)
}

func FromContext(ctx context.Context) manta.Authorizer {
	v := ctx.Value(authorizerCtxKey)
	if v == nil {
		return nil
	}

	return v.(manta.Authorizer)
}
