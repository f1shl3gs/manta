package all

import (
	"context"

	"github.com/f1shl3gs/manta/kv"
)

func Migration0001AddSecret() Spec {
	bucket := []byte("secrets")

	return &spec{
		name: "add secret",
		up: func(ctx context.Context, store kv.SchemaStore) error {
			return store.CreateBucket(ctx, bucket)
		},
	}
}
