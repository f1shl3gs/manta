package kv_test

import (
	"context"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/mock"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

const (
	mockUID = manta.ID(1)
)

type initSessionService func(t *testing.T) (context.Context, manta.SessionService, func())

func TestSessionService(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T, init initSessionService)
	}{
		{
			name: "CreateSession",
			fn:   CreatSession,
		},
		{
			name: "FindSession",
			fn:   FindSession,
		},
		{
			name: "RenewSession",
			fn:   RenewSession,
		},
		{
			name: "RevokeSession",
			fn:   RevokeSession,
		},
	}

	var init initSessionService = func(t *testing.T) (context.Context, manta.SessionService, func()) {
		idGen := mock.NewIncrementalIDGenerator(1)
		svc, closer := NewTestService(t, kv.WithIDGenerator(idGen))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)

		// create a mock user
		err := svc.CreateUser(ctx, &manta.User{
			Name: "mock",
		})
		require.NoError(t, err)

		return ctx, svc, func() {
			closer()
			cancel()
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			tt.fn(t, init)
		})
	}
}

func CreatSession(t *testing.T, init initSessionService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.SessionService)
	}{
		{
			name: "create with invalid uid",
			fn: func(t *testing.T, ctx context.Context, svc manta.SessionService) {
				uid := manta.ID(0)
				_, err := svc.CreateSession(ctx, uid)
				require.Equal(t, manta.ErrInvalidID, err)
			},
		},
		{
			name: "create multiple times",
			fn: func(t *testing.T, ctx context.Context, svc manta.SessionService) {
				uid := manta.ID(1)
				s1, err := svc.CreateSession(ctx, uid)
				require.NoError(t, err)
				require.Equal(t, s1.UserID, uid)
				require.EqualValues(t, s1.ID, 2)

				s2, err := svc.CreateSession(ctx, uid)
				require.NoError(t, err)
				require.Equal(t, s2.UserID, uid)
				require.EqualValues(t, s2.ID, 3)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, svc, closer := init(t)
			defer closer()

			tt.fn(t, ctx, svc)
		})
	}
}

func FindSession(t *testing.T, init initSessionService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.SessionService)
	}{
		{
			name: "find",
			fn: func(t *testing.T, ctx context.Context, svc manta.SessionService) {
				session, err := svc.CreateSession(ctx, mockUID)
				require.NoError(t, err)
				require.Equal(t, mockUID, session.UserID)

				found, err := svc.FindSession(ctx, session.ID)
				require.NoError(t, err)

				diff := cmp.Diff(found, session)
				require.Equal(t, "", diff, "created and found is not equal\n%s", diff)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, svc, closer := init(t)
			defer closer()

			tt.fn(t, ctx, svc)
		})
	}
}

func RenewSession(t *testing.T, init initSessionService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.SessionService)
	}{
		{
			name: "renew",
			fn: func(t *testing.T, ctx context.Context, svc manta.SessionService) {
				session, err := svc.CreateSession(ctx, mockUID)
				require.NoError(t, err)
				require.Equal(t, mockUID, session.UserID)

				expiration := time.Now().Add(time.Hour)
				err = svc.RenewSession(ctx, session.ID, expiration)
				require.NoError(t, err)

				ns, err := svc.FindSession(ctx, session.ID)
				require.NoError(t, err)

				require.Equal(t, true, ns.ExpiresAt.UnixNano() == expiration.UnixNano())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, svc, closer := init(t)
			defer closer()

			tt.fn(t, ctx, svc)
		})
	}
}

func RevokeSession(t *testing.T, init initSessionService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.SessionService)
	}{
		{
			name: "revoke",
			fn: func(t *testing.T, ctx context.Context, svc manta.SessionService) {
				session, err := svc.CreateSession(ctx, mockUID)
				require.NoError(t, err)
				require.Equal(t, mockUID, session.UserID)

				err = svc.RevokeSession(ctx, session.ID)
				require.NoError(t, err)
			},
		},
		{
			name: "revoke none exist",
			fn: func(t *testing.T, ctx context.Context, svc manta.SessionService) {
				err := svc.RevokeSession(ctx, manta.ID(2))
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, svc, closer := init(t)
			defer closer()

			tt.fn(t, ctx, svc)
		})
	}
}
