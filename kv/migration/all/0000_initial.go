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
			[]byte("buckets"),
			[]byte("bucketorgindex"),
			[]byte("checks"),
			[]byte("checkorgindex"),
			[]byte("dashboards"),
			[]byte("dashboardorgindex"),
			[]byte("events"),
			[]byte("eventorgindex"),
			[]byte("keyring"),
			[]byte("kvlog"),
			[]byte("notifications"),
			[]byte("notificationorgindex"),
			[]byte("notificationendpoints"),
			[]byte("notificationendpointorgindex"),
			[]byte("organizations"),
			[]byte("organizationnameindex"),
			[]byte("otcls"),
			[]byte("otclorgindex"),
			[]byte("passwords"),
			[]byte("scrapertargets"),
			[]byte("scrapertargetorgindex"),
			[]byte("sessions"),
			[]byte("tasks"),
			[]byte("taskorgindex"),
			[]byte("taskownerindex"),
			[]byte("runs"),
			[]byte("runtaskindex"),
			[]byte("templates"),
			[]byte("templatenameindex"),
			[]byte("users"),
			[]byte("usernameindex"),
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
	}
}
