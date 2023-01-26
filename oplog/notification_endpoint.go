package oplog

import (
	"context"
	"encoding/json"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"

	"go.uber.org/zap"
)

type NotificationEndpointService struct {
	manta.NotificationEndpointService

	logger *zap.Logger
	oplog  manta.OperationLogService
}

func NewNotificationEndpointService(
	service manta.NotificationEndpointService,
	oplog manta.OperationLogService,
	logger *zap.Logger,
) *NotificationEndpointService {
	return &NotificationEndpointService{
		NotificationEndpointService: service,
		logger:                      logger,
		oplog:                       oplog,
	}
}

// CreateNotificationEndpoint creates a new notification endpoint and sets b.ID with the new identifier
func (s *NotificationEndpointService) CreateNotificationEndpoint(
	ctx context.Context,
	ne manta.NotificationEndpoint,
) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.NotificationEndpointService.CreateNotificationEndpoint(ctx, ne)
	if err != nil {
		return err
	}

	ne, err = s.NotificationEndpointService.FindNotificationEndpointByID(ctx, ne.GetID())
	if err != nil {
		return err
	}

	data, err := json.Marshal(ne)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Create,
		ResourceID:   ne.GetID(),
		ResourceType: manta.NotificationEndpointsResourceType,
		OrgID:        ne.GetOrgID(),
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add create notification endpoint oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", ne.GetID()),
			zap.Stringer("orgID", ne.GetOrgID()))
	}

	return err
}

// UpdateNotificationEndpoint updates a single notification endpoint.
// Returns the new notification endpoint after update.
func (s *NotificationEndpointService) UpdateNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
	ne manta.NotificationEndpoint,
) (manta.NotificationEndpoint, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	ne, err = s.NotificationEndpointService.UpdateNotificationEndpoint(ctx, id, ne)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(ne)
	if err != nil {
		return nil, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   ne.GetID(),
		ResourceType: manta.NotificationEndpointsResourceType,
		OrgID:        ne.GetOrgID(),
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add update notification endpoint oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", ne.GetID()),
			zap.Stringer("orgID", ne.GetOrgID()))
		return nil, err
	}

	return ne, nil
}

// PatchNotificationEndpoint patch a single notification endpoint.
// Returns the new notification endpoint after patch.
func (s *NotificationEndpointService) PatchNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
	upd manta.NotificationEndpointUpdate,
) (manta.NotificationEndpoint, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	ne, err := s.NotificationEndpointService.PatchNotificationEndpoint(ctx, id, upd)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(ne)
	if err != nil {
		return nil, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   ne.GetID(),
		ResourceType: manta.NotificationEndpointsResourceType,
		OrgID:        ne.GetOrgID(),
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add patch notification endpoint failed",
			zap.Error(err),
			zap.Stringer("resourceID", ne.GetID()),
			zap.Stringer("orgID", ne.GetOrgID()))
		return nil, err
	}

	return ne, nil
}

// DeleteNotificationEndpoint remove a notification endpoint by ID, return it's secret
// fields, orgID for further deletion
func (s *NotificationEndpointService) DeleteNotificationEndpoint(
	ctx context.Context,
	id manta.ID,
) ([]manta.SecretField, manta.ID, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, 0, err
	}

	now := time.Now()
	sfs, orgID, err := s.NotificationEndpointService.DeleteNotificationEndpoint(ctx, id)
	if err != nil {
		return nil, 0, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Delete,
		ResourceID:   id,
		ResourceType: manta.NotificationEndpointsResourceType,
		OrgID:        orgID,
		UserID:       auth.GetUserID(),
		ResourceBody: nil,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add delete notification endpoint oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", id),
			zap.Stringer("orgID", orgID))
		return nil, 0, err
	}

	return sfs, orgID, nil
}
