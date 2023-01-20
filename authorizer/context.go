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

func FromContext(ctx context.Context) (manta.Authorizer, error) {
	a, ok := ctx.Value(authorizerCtxKey).(manta.Authorizer)
	if !ok {
        return nil, &manta.Error{
            Code: manta.EInternal,
			Msg:  "authorizer not found on context",
		}
	}

	if a == nil {
        return nil, &manta.Error{
            Code: manta.EInternal,
			Msg:  "unexpected authorizer",
		}
	}

	return a, nil
}
