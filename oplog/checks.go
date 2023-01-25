package oplog

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"
)

type CheckService struct {
	manta.CheckService

	logger *zap.Logger
	oplog  manta.OperationLogService
}

func NewCheckService(checkService manta.CheckService, oplog manta.OperationLogService, logger *zap.Logger) *CheckService {
	return &CheckService{
		CheckService: checkService,
		oplog:        oplog,
	}
}

// CreateCheck creates a new and set its id with new identifier
func (s *CheckService) CreateCheck(ctx context.Context, c *manta.Check) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.CheckService.CreateCheck(ctx, c)
	if err != nil {
		return err
	}

	c, err = s.CheckService.FindCheckByID(ctx, c.ID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Create,
		ResourceID:   c.ID,
		ResourceType: manta.ChecksResourceType,
		OrgID:        c.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add create check oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", c.ID),
			zap.Stringer("orgID", c.OrgID))
	}

	return nil
}

// UpdateCheck updates the whole check, returns the new check after update
func (s *CheckService) UpdateCheck(ctx context.Context, id manta.ID, c *manta.Check) (*manta.Check, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	c, err = s.CheckService.UpdateCheck(ctx, id, c)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   id,
		ResourceType: manta.ChecksResourceType,
		OrgID:        c.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add update check oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", c.ID),
			zap.Stringer("orgID", c.OrgID))
	}

	return c, nil
}

// PatchCheck updates a single check with changeset
// Returns the new check after patch
func (s *CheckService) PatchCheck(ctx context.Context, id manta.ID, upd manta.CheckUpdate) (*manta.Check, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	c, err := s.CheckService.PatchCheck(ctx, id, upd)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   id,
		ResourceType: manta.ChecksResourceType,
		OrgID:        c.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add patch check log failed",
			zap.Error(err),
			zap.Stringer("resourceID", c.ID),
			zap.Stringer("orgID", c.OrgID))
	}

	return c, nil
}

// DeleteCheck delete a single check by ID
func (s *CheckService) DeleteCheck(ctx context.Context, id manta.ID) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	c, err := s.CheckService.FindCheckByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.CheckService.DeleteCheck(ctx, id)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Delete,
		ResourceID:   id,
		ResourceType: manta.ChecksResourceType,
		OrgID:        c.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: nil,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add delete check log failed",
			zap.Error(err),
			zap.Stringer("resourceID", c.ID),
			zap.Stringer("orgID", c.OrgID))
	}

	return nil
}
