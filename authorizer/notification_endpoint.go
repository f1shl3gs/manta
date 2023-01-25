package authorizer

import (
	"context"

	"github.com/f1shl3gs/manta"
)

type NotificationEndpointService struct {
	// service manta.NotificationEndpointService
}

// FindNotificationEndpointByID returns a single notification endpoint by ID
func (s *NotificationEndpointService) FindNotificationEndpointByID(
	ctx context.Context,
	id manta.ID,
) (manta.NotificationEndpoint, error) {
	panic("not implement")
}

// FindNotificationEndpoints returns a list of notication endpoints that match filter.
func (s *NotificationEndpointService) FindNotificationEndpoints(
	ctx context.Context,
	filter manta.NotificationEndpointFilter,
) ([]manta.NotificationEndpoint, error) {
	panic("not implement")
}

// CreateNotificationEndpoint creates a new notification endpoint and sets b.ID with the new identifier
func (s *NotificationEndpointService) CreateNotificationEndpoint(
	ctx context.Context,
	ne manta.NotificationEndpoint,
) error {
	panic("not implement")
}

// UpdateNotificationEndpoint updates a single notification endpoint.
// Returns the new notification endpoint after update.
func (s *NotificationEndpointService) UpdateNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
	ne manta.NotificationEndpoint,
) (manta.NotificationEndpoint, error) {
	panic("not implement")
}

// PatchNotificationEndpoint patch a single notification endpoint.
// Returns the new notification endpoint after patch.
func (s *NotificationEndpointService) PatchNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
	upd manta.NotificationEndpointUpdate,
) (manta.NotificationEndpoint, error) {
	panic("not implement")
}

// DeleteNotificationEndpoint remove a notification endpoint by ID, return it's secret fields,
// orgID for further deletion
func (s *NotificationEndpointService) DeleteNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
) ([]manta.SecretField, manta.ID, error) {
	panic("not implement")
}
