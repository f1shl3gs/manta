package all

import (
	"context"

	"github.com/f1shl3gs/manta/kv"
)

func Migration0000Initial() Spec {
	var (
		buckets = [][]byte{
			[]byte("authorizations"),
			[]byte("authorizationtokenindex"),
			[]byte("authorizationuserindex"),
			kv.ChecksBucket,
			kv.CheckOrgIndexBucket,
			kv.ConfigurationBucket,
			kv.ConfigurationOrgIndexBucket,
			kv.DashboardsBucket,
			kv.DashboardOrgIndexBucket,
			kv.CellsBucket,
			kv.CellDashboardIndexBucket,
			[]byte("organizations"),
			[]byte("organizationnameindex"),
			[]byte("passwords"),
			[]byte("sessions"),
			[]byte("users"),
			[]byte("usernameindex"),
			kv.RegistryBucket,
			kv.ScraperBucket,
			kv.ScrapeOrgIndexBucket,
			kv.TasksBucket,
			kv.TaskOrgIndexBucket,
			kv.TaskOwnerIndexBucket,
		}
	)

	return &spec{
		name: "initial",
		up: func(ctx context.Context, store kv.SchemaStore) error {
			for _, b := range buckets {
				err := store.CreateBucket(ctx, b)
				if err != nil {
					return err
				}
			}

			return nil
		},
		down: func(ctx context.Context, store kv.SchemaStore) error {
			panic("initial store cannot be downgrade")
		},
	}
}
