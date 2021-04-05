package all

import (
	"context"

	"github.com/f1shl3gs/manta/kv"
)

func Migration0001AddSecret() Spec {
	buckets := [][]byte{
		[]byte("secrets"),
		[]byte("secretorgindex"),
	}

	return &spec{
		name: "add secret",
		up: func(ctx context.Context, store kv.SchemaStore) error {
			for _, b := range buckets {
				err := store.CreateBucket(ctx, b)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
