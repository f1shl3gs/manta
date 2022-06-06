package kv_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/kv"
)

func newTestPasswordService(t *testing.T, timeout time.Duration) (context.Context, manta.PasswordService, func()) {
	svc, closer := NewTestService(t)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	return ctx, svc, func() {
		closer()
		cancel()
	}
}

func TestPasswordService(t *testing.T) {
	const (
		uid      = manta.ID(1)
		password = "123456789"
	)

	t.Run("set and compare", func(t *testing.T) {
		ctx, svc, closer := newTestPasswordService(t, 2*time.Second)
		defer closer()

		err := svc.SetPassword(ctx, uid, password)
		require.NoError(t, err)

		err = svc.ComparePassword(ctx, uid, password)
		require.NoError(t, err)

		err = svc.ComparePassword(ctx, uid, "aaaa")
		require.Equal(t, manta.ErrPasswordNotMatch, err)
	})

	t.Run("delete password", func(t *testing.T) {
		ctx, svc, closer := newTestPasswordService(t, 2*time.Second)
		defer closer()

		err := svc.SetPassword(ctx, uid, password)
		require.NoError(t, err)

		err = svc.ComparePassword(ctx, uid, password)
		require.NoError(t, err)

		err = svc.DeletePassword(ctx, uid)
		require.NoError(t, err)

		err = svc.ComparePassword(ctx, uid, password)
		require.Equal(t, kv.ErrKeyNotFound, err)

		err = svc.DeletePassword(ctx, uid)
		require.NoError(t, err)
	})
}
