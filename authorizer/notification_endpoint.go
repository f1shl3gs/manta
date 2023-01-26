package authorizer

import (
	"context"

	"github.com/f1shl3gs/manta"
)

type NotificationEndpointService struct {
	service manta.NotificationEndpointService
}

var _ manta.NotificationEndpointService = &NotificationEndpointService{}

func NewNotificationEndpointService(service manta.NotificationEndpointService) *NotificationEndpointService {
	return &NotificationEndpointService{
		service: service,
	}
}

// FindNotificationEndpointByID returns a single notification endpoint by ID
func (s *NotificationEndpointService) FindNotificationEndpointByID(
	ctx context.Context,
	id manta.ID,
) (manta.NotificationEndpoint, error) {
	ne, err := s.service.FindNotificationEndpointByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeRead(ctx, manta.NotificationEndpointsResourceType, id, ne.GetOrgID()); err != nil {
		return nil, err
	}

	return ne, nil
}

// FindNotificationEndpoints returns a list of notication endpoints that match filter.
func (s *NotificationEndpointService) FindNotificationEndpoints(
	ctx context.Context,
	filter manta.NotificationEndpointFilter,
) ([]manta.NotificationEndpoint, error) {
	nes, err := s.service.FindNotificationEndpoints(ctx, filter)
	if err != nil {
		return nil, err
	}

	filtered := nes[:0]
	for _, ne := range nes {
		_, _, err := authorizeRead(ctx, manta.NotificationEndpointsResourceType, ne.GetID(), ne.GetOrgID())
		if err != nil && manta.ErrorCode(err) != manta.EUnauthorized {
			return nil, err
		}

		if manta.ErrorCode(err) == manta.EUnauthorized {
			continue
		}

		filtered = append(filtered, ne)
	}

	return nes, nil
}

// CreateNotificationEndpoint creates a new notification endpoint and sets b.ID with the new identifier
func (s *NotificationEndpointService) CreateNotificationEndpoint(
	ctx context.Context,
	ne manta.NotificationEndpoint,
) error {
	if _, _, err := authorizeOrgWriteResource(ctx, manta.NotificationEndpointsResourceType, ne.GetOrgID()); err != nil {
		return err
	}

	return s.service.CreateNotificationEndpoint(ctx, ne)
}

// UpdateNotificationEndpoint updates a single notification endpoint.
// Returns the new notification endpoint after update.
func (s *NotificationEndpointService) UpdateNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
	ne manta.NotificationEndpoint,
) (manta.NotificationEndpoint, error) {
	curr, err := s.service.FindNotificationEndpointByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeWrite(ctx, manta.NotificationEndpointsResourceType, id, curr.GetOrgID()); err != nil {
		return nil, err
	}

	return s.service.UpdateNotificationEndpoint(ctx, id, ne)
}

// PatchNotificationEndpoint patch a single notification endpoint.
// Returns the new notification endpoint after patch.
func (s *NotificationEndpointService) PatchNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
	upd manta.NotificationEndpointUpdate,
) (manta.NotificationEndpoint, error) {
	ne, err := s.service.FindNotificationEndpointByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeWrite(ctx, manta.NotificationEndpointsResourceType, id, ne.GetOrgID()); err != nil {
		return nil, err
	}

	return s.service.PatchNotificationEndpoint(ctx, id, upd)
}

// DeleteNotificationEndpoint remove a notification endpoint by ID, return it's secret fields,
// orgID for further deletion
func (s *NotificationEndpointService) DeleteNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
) ([]manta.SecretField, manta.ID, error) {
	ne, err := s.service.FindNotificationEndpointByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}

	if _, _, err := authorizeWrite(ctx, manta.NotificationEndpointsResourceType, id, ne.GetOrgID()); err != nil {
		return nil, 0, err
	}

	return s.service.DeleteNotificationEndpoint(ctx, id)
}
