package kv_test

import (
	"context"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/mock"

	"github.com/stretchr/testify/assert"
)

var (
	testResourceID_1 = manta.ID(1)
	testResourceID_2 = manta.ID(2)
	testUserID_1     = manta.ID(1)
	testUserID_2     = manta.ID(2)
)

type initOpLogService func(t *testing.T) (context.Context, manta.OperationLogService, func())

func TestOpLogService(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T, init initOpLogService)
	}{
		{
			name: "AddLogEntry",
			fn:   AddLogEntry,
		},
		{
			name: "FindOperationLogsByID",
			fn:   FindOperationLogsByID,
		},
		{
			name: "FindOperationLogsByUser",
			fn:   FindOperationLogsByUser,
		},
	}

	var initFn initOpLogService = func(t *testing.T) (context.Context, manta.OperationLogService, func()) {
		idGen := mock.NewIncrementalIDGenerator(1)
		svc, closer := NewTestService(t, kv.WithIDGenerator(idGen))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err := svc.AddLogEntry(ctx, manta.OperationLogEntry{
			Type:         manta.Create,
			ResourceID:   testResourceID_1,
			ResourceType: manta.ChecksResourceType,
			OrgID:        1,
			UserID:       testUserID_2,
			ResourceBody: nil,
			Time:         time.Now(),
		})
		assert.NoError(t, err)

		err = svc.AddLogEntry(ctx, manta.OperationLogEntry{
			Type:         manta.Create,
			ResourceID:   testResourceID_2,
			ResourceType: manta.DashboardsResourceType,
			OrgID:        1,
			UserID:       testUserID_2,
			ResourceBody: nil,
			Time:         time.Now(),
		})
		assert.NoError(t, err)

		err = svc.AddLogEntry(ctx, manta.OperationLogEntry{
			Type:         manta.Update,
			ResourceID:   testResourceID_2,
			ResourceType: manta.DashboardsResourceType,
			OrgID:        1,
			UserID:       testUserID_1,
			ResourceBody: nil,
			Time:         time.Now(),
		})
		assert.NoError(t, err)

		return ctx, svc, func() {
			cancel()
			closer()
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

func AddLogEntry(t *testing.T, initFn initOpLogService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.OperationLogService)
	}{
		{
			name: "create one",
			fn: func(t *testing.T, ctx context.Context, svc manta.OperationLogService) {
				svc.AddLogEntry(ctx, manta.OperationLogEntry{
					Type:         manta.Update,
					ResourceID:   1,
					ResourceType: manta.ChecksResourceType,
					OrgID:        1,
					UserID:       1,
					ResourceBody: nil,
					Time:         time.Now(),
				})
			},
		},
		{
			name: "add multiple times",
			fn: func(t *testing.T, ctx context.Context, svc manta.OperationLogService) {
				c := manta.OperationLogEntry{
					Type:         manta.Update,
					ResourceID:   1,
					ResourceType: manta.ChecksResourceType,
					OrgID:        1,
					UserID:       1,
					ResourceBody: nil,
					Time:         time.Now(),
				}

				svc.AddLogEntry(ctx, c)
				time.Sleep(time.Second)
				c.Time = time.Now()
				svc.AddLogEntry(ctx, c)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, svc, closer := initFn(t)
			defer closer()

			tt.fn(t, ctx, svc)
		})
	}
}

func FindOperationLogsByID(t *testing.T, initFn initOpLogService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.OperationLogService)
	}{
		{
			name: "find",
			fn: func(t *testing.T, ctx context.Context, svc manta.OperationLogService) {
				changes, _, err := svc.FindOperationLogsByID(ctx, testResourceID_1, manta.FindOptions{
					Limit: 10,
				})
				assert.NoError(t, err)
				assert.Equal(t, 1, len(changes))

				changes, _, err = svc.FindOperationLogsByID(ctx, testResourceID_2, manta.FindOptions{
					Limit: 10,
				})
				assert.NoError(t, err)
				assert.Equal(t, 2, len(changes))
			},
		},
		{
			name: "find with limit",
			fn: func(t *testing.T, ctx context.Context, svc manta.OperationLogService) {
				changes, _, err := svc.FindOperationLogsByID(ctx, testResourceID_2, manta.FindOptions{
					Limit: 1,
				})
				assert.NoError(t, err)
				assert.Equal(t, 1, len(changes))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, svc, closer := initFn(t)
			defer closer()

			tt.fn(t, ctx, svc)
		})
	}
}

func FindOperationLogsByUser(t *testing.T, initFn initOpLogService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.OperationLogService)
	}{
		{
			name: "find",
			fn: func(t *testing.T, ctx context.Context, svc manta.OperationLogService) {
				changes, _, err := svc.FindOperationLogsByUser(ctx, testUserID_1, manta.FindOptions{
					Limit: 10,
				})
				assert.NoError(t, err)
				assert.Equal(t, 1, len(changes))

				changes, _, err = svc.FindOperationLogsByUser(ctx, testUserID_2, manta.FindOptions{
					Limit: 10,
				})
				assert.NoError(t, err)
				assert.Equal(t, 2, len(changes))
			},
		},
		{
			name: "find with limit",
			fn: func(t *testing.T, ctx context.Context, svc manta.OperationLogService) {
				changes, _, err := svc.FindOperationLogsByUser(ctx, testUserID_2, manta.FindOptions{
					Limit: 1,
				})
				assert.NoError(t, err)
				assert.Equal(t, 1, len(changes))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, svc, closer := initFn(t)
			defer closer()

			tt.fn(t, ctx, svc)
		})
	}
}
