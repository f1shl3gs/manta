package kv_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestKeyring(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		svc, closer := NewTestService(t)
		defer closer()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := svc.AddKey(ctx, []byte("foo"))
		require.NoError(t, err)

		pk, err := svc.PrimaryKey(ctx)
		require.NoError(t, err)

		require.Equal(t, []byte("foo"), pk)
	})

	t.Run("remove primary key", func(t *testing.T) {
		svc, closer := NewTestService(t)
		defer closer()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		input := [][]byte{
			[]byte("foo"),
			[]byte("bar"),
		}

		for _, key := range input {
			err := svc.AddKey(ctx, key)
			require.NoError(t, err)
		}

		pk, err := svc.PrimaryKey(ctx)
		require.NoError(t, err)
		require.Equal(t, []byte("bar"), pk)

		err = svc.RemoveKey(ctx, pk)
		require.NoError(t, err)

		pk, err = svc.PrimaryKey(ctx)
		require.NoError(t, err)
		require.Equal(t, []byte("foo"), pk)
	})
}
