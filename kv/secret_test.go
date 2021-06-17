package kv_test

import (
	"context"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
)

type initSecretService func(t *testing.T) (context.Context, manta.SessionService, func())

func TestSecretService(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T, init initSecretService)
	}{
		{
			name: "GetSecretKeys",
			fn:   GetSecretKeys,
		},
	}

	var initFn initSecretService = func(t *testing.T) (context.Context, manta.SessionService, func()) {
		svc, closer := NewTestService(t)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		return ctx, svc, func() {
			closer()
			cancel()
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			tt.fn(t, initFn)
		})
	}
}

func GetSecretKeys(t *testing.T, init initSecretService) {

}
