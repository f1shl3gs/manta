package all

import (
	"context"

	"github.com/f1shl3gs/manta/kv"
)

func Migration0002AddVariable() Spec {
	buckets := [][]byte{
		[]byte("variables"),
		[]byte("variableorgindex"),
	}

	return &spec{
		name: "add variable",
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
