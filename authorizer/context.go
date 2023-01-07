package authorizer

import (
	"context"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/errors"
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
		return nil, &errors.Error{
			Code: errors.EInternal,
			Msg:  "authorizer not found on context",
		}
	}

	if a == nil {
		return nil, &errors.Error{
			Code: errors.EInternal,
			Msg:  "unexpected authorizer",
		}
	}

	return a, nil
}
