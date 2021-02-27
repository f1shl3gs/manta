package manta

import (
	"context"
)

type NotificationEndpointFilter struct {
	ID    *ID
	Name  *string
	OrgID *ID
}

type NotificationEndpointUpdate struct {
	Name        *string
	Description *string
}

type NotificationEndpointService interface {
	// FindNotificationEndpointByID find notification endpoint by id
	FindNotificationEndpointByID(ctx context.Context, id ID) (*NotificationEndpoint, error)

	// FindNotificationEndpoints returns a list of notification endpoints that match filter and the total count of matching notification endpoints
	// additional options provide pagination & sorting
	FindNotificationEndpoints(ctx context.Context, filter NotificationEndpointFilter, opt ...FindOptions) ([]*NotificationEndpoint, int, error)

	// CreateNotificationEndpoint create a single notification endpoint and sets id with the new identifier
	CreateNotificationEndpoint(ctx context.Context, ne *NotificationEndpoint) error

	// UpdateNotificationEndpoint update a single notification endpoint
	UpdateNotificationEndpoint(ctx context.Context, id ID, u NotificationEndpointUpdate) (*NotificationEndpoint, error)

	// DeleteNotificationEndpoint delete a notification endpoint by id
	DeleteNotificationEndpoint(ctx context.Context, id ID) error
}
