package kv_test

import (
	"context"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/kv"
	"github.com/f1shl3gs/manta/mock"
	"github.com/f1shl3gs/manta/notification"

	"github.com/stretchr/testify/assert"
)

var ()

type initNotificationEndpointService func(t *testing.T) (context.Context, manta.NotificationEndpointService, func())

func TestNotificationEndpointService(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T, initFn initNotificationEndpointService)
	}{
		{
			name: "FindNotificationEndpointByID",
			fn:   FindNotificationEndpointByID,
		},
		{
			name: "FindNotificationEndpoints",
			fn:   FindNotificationEndpoints,
		},
		{
			name: "CreateNotificationEndpoint",
			fn:   CreateNotificationEndpoint,
		},
		{
			name: "PatchNotificationEndpoint",
			fn:   PatchNotificationEndpoint,
		},
		{
			name: "DeleteNotificationEndpoint",
			fn:   DeleteNotificationEndpoint,
		},
	}

	var initFn initNotificationEndpointService = func(t *testing.T) (context.Context, manta.NotificationEndpointService, func()) {
		idGen := mock.NewIncrementalIDGenerator(1)
		svc, closer := NewTestService(t, kv.WithIDGenerator(idGen))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		ne1 := &notification.HTTP{
			Base: notification.Base{
				ID:      0,
				Created: time.Time{},
				Updated: time.Time{},
				Name:    "foo",
				Desc:    "foo desc",
				OrgID:   1,
			},
			URL:             "http://example.com",
			Headers:         nil,
			Username:        manta.SecretField{},
			Password:        manta.SecretField{},
			Token:           manta.SecretField{},
			Method:          "POST",
			AuthMethod:      "none",
			ContentTemplate: "",
		}

		ne2 := &notification.HTTP{
			Base: notification.Base{
				ID:      0,
				Created: time.Time{},
				Updated: time.Time{},
				Name:    "bar",
				Desc:    "bar desc",
				OrgID:   2,
			},
			URL:             "http://example.com",
			Headers:         nil,
			Username:        manta.SecretField{},
			Password:        manta.SecretField{},
			Token:           manta.SecretField{},
			Method:          "POST",
			AuthMethod:      "none",
			ContentTemplate: "",
		}

		err := svc.CreateNotificationEndpoint(ctx, ne1)
		assert.NoError(t, err)
		err = svc.CreateNotificationEndpoint(ctx, ne2)
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

func FindNotificationEndpointByID(t *testing.T, initFn initNotificationEndpointService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService)
	}{
		{
			name: "find",
			fn: func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService) {
				ne, err := svc.FindNotificationEndpointByID(ctx, manta.ID(1))
				assert.NoError(t, err)
				assert.Equal(t, manta.ID(1), ne.GetID())

				ne, err = svc.FindNotificationEndpointByID(ctx, manta.ID(2))
				assert.NoError(t, err)
				assert.Equal(t, manta.ID(2), ne.GetID())
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

func FindNotificationEndpoints(t *testing.T, initFn initNotificationEndpointService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService)
	}{
		{
			name: "find",
			fn: func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService) {
				nes, err := svc.FindNotificationEndpoints(ctx, manta.NotificationEndpointFilter{
					OrgID: 1,
				})
				assert.NoError(t, err)
				assert.Equal(t, 1, len(nes))

				nes, err = svc.FindNotificationEndpoints(ctx, manta.NotificationEndpointFilter{
					OrgID: 2,
				})
				assert.NoError(t, err)
				assert.Equal(t, 1, len(nes))
			},
		},
		{
			name: "find with unknown orgID",
			fn: func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService) {
				nes, err := svc.FindNotificationEndpoints(ctx, manta.NotificationEndpointFilter{
					OrgID: 3,
				})
				assert.NoError(t, err)
				assert.Equal(t, 0, len(nes))
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

func CreateNotificationEndpoint(t *testing.T, initFn initNotificationEndpointService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService)
	}{
		{
			name: "create one",
			fn: func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService) {
				var ne manta.NotificationEndpoint = &notification.HTTP{
					Base: notification.Base{
						ID:      0,
						Created: time.Time{},
						Updated: time.Time{},
						Name:    "foo",
						Desc:    "foo desc",
						OrgID:   1,
					},
					URL:             "http://example.com",
					Headers:         nil,
					Username:        manta.SecretField{},
					Password:        manta.SecretField{},
					Token:           manta.SecretField{},
					Method:          "POST",
					AuthMethod:      "none",
					ContentTemplate: "",
				}

				err := svc.CreateNotificationEndpoint(ctx, ne)
				assert.NoError(t, err)

				ne, err = svc.FindNotificationEndpointByID(ctx, manta.ID(3))
				assert.NoError(t, err)
				assert.Equal(t, manta.ID(3), ne.GetID())

				nes, err := svc.FindNotificationEndpoints(ctx, manta.NotificationEndpointFilter{
					OrgID: 1,
				})
				assert.NoError(t, err)
				assert.Equal(t, 2, len(nes))
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

func UpdateNotificationEndpoint(t *testing.T, initFn initNotificationEndpointService) {
}
func PatchNotificationEndpoint(t *testing.T, initFn initNotificationEndpointService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService)
	}{
		{
			name: "update name",
			fn: func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService) {
				newName := "fooo"
				newDesc := "descccc"

				ne, err := svc.PatchNotificationEndpoint(ctx, 1, manta.NotificationEndpointUpdate{
					Name: &newName,
					Desc: &newDesc,
				})
				assert.NoError(t, err)
				h := ne.(*notification.HTTP)
				assert.Equal(t, newName, h.Name)
				assert.Equal(t, newDesc, h.Desc)
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

func DeleteNotificationEndpoint(t *testing.T, initFn initNotificationEndpointService) {
	tests := []struct {
		name string
		fn   func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService)
	}{
		{
			name: "delete",
			fn: func(t *testing.T, ctx context.Context, svc manta.NotificationEndpointService) {
				sfs, orgID, err := svc.DeleteNotificationEndpoint(ctx, 1)
				assert.NoError(t, err)
				assert.Equal(t, 0, len(sfs))
				assert.Equal(t, manta.ID(1), orgID)

				ne, err := svc.FindNotificationEndpointByID(ctx, 1)
				assert.Nil(t, ne)
				assert.Equal(t, kv.ErrKeyNotFound, err)
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
