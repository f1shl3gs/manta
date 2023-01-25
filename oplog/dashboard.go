package oplog

import (
	"context"
	"encoding/json"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authorizer"

	"go.uber.org/zap"
)

type DashboardService struct {
	manta.DashboardService

	logger *zap.Logger
	oplog  manta.OperationLogService
}

var _ manta.DashboardService = &DashboardService{}

func NewDashboardService(dashboardService manta.DashboardService, oplog manta.OperationLogService, logger *zap.Logger) *DashboardService {
	return &DashboardService{
		DashboardService: dashboardService,
		logger:           logger,
		oplog:            oplog,
	}
}

func (s *DashboardService) CreateDashboard(ctx context.Context, d *manta.Dashboard) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.DashboardService.CreateDashboard(ctx, d)
	if err != nil {
		return err
	}

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Create,
		ResourceID:   d.ID,
		ResourceType: manta.DashboardsResourceType,
		OrgID:        d.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add create dashbaord oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", d.ID),
			zap.Stringer("orgID", d.OrgID))
	}

	return err
}

func (s *DashboardService) UpdateDashboard(ctx context.Context, upd manta.DashboardUpdate) (*manta.Dashboard, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dashboard, err := s.DashboardService.UpdateDashboard(ctx, upd)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(dashboard)
	if err != nil {
		return nil, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   dashboard.ID,
		ResourceType: manta.DashboardsResourceType,
		OrgID:        dashboard.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add update dashbaord oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", dashboard.ID),
			zap.Stringer("orgID", dashboard.OrgID))
		return nil, err
	}

	return dashboard, nil
}

func (s *DashboardService) AddDashboardCell(ctx context.Context, dashboardID manta.ID, cell *manta.Cell) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.DashboardService.AddDashboardCell(ctx, dashboardID, cell)
	if err != nil {
		return err
	}

	// post part
	dashboard, err := s.DashboardService.FindDashboardByID(ctx, dashboardID)
	if err != nil {
		s.logger.Error("find dashboard to add cell oplog failed",
			zap.Error(err),
			zap.Stringer("dashboard", dashboardID))
		return err
	}

	data, err := json.Marshal(dashboard)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   dashboardID,
		ResourceType: manta.DashboardsResourceType,
		OrgID:        dashboard.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add create cell oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", dashboard.ID),
			zap.Stringer("orgID", dashboard.OrgID))
		return err
	}

	return nil
}

func (s *DashboardService) RemoveDashboardCell(ctx context.Context, dashboardID, cellID manta.ID) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.DashboardService.RemoveDashboardCell(ctx, dashboardID, cellID)
	if err != nil {
		return err
	}

	// TODO: stale data?
	dashboard, err := s.DashboardService.FindDashboardByID(ctx, dashboardID)
	if err != nil {
		s.logger.Error("find dashboard to add cell oplog failed",
			zap.Error(err),
			zap.Stringer("dashboard", dashboardID))
		return err
	}

	data, err := json.Marshal(dashboard)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   dashboardID,
		ResourceType: manta.DashboardsResourceType,
		OrgID:        dashboard.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add remove cell oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", dashboard.ID),
			zap.Stringer("orgID", dashboard.OrgID))
		return err
	}

	return nil
}

func (s *DashboardService) UpdateDashboardCell(ctx context.Context, upd manta.DashboardCellUpdate) (*manta.Cell, error) {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cell, err := s.DashboardService.UpdateDashboardCell(ctx, upd)
	if err != nil {
		return nil, err
	}

	dashboard, err := s.DashboardService.FindDashboardByID(ctx, upd.DashboardID)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(dashboard)
	if err != nil {
		return nil, err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   dashboard.ID,
		ResourceType: manta.DashboardsResourceType,
		OrgID:        dashboard.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add update cell oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", dashboard.ID),
			zap.Stringer("orgID", dashboard.OrgID))
		return nil, err
	}

	return cell, nil
}

func (s *DashboardService) DeleteDashboard(ctx context.Context, id manta.ID) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()

	dashboard, err := s.DashboardService.FindDashboardByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.DashboardService.RemoveDashboard(ctx, id)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   dashboard.ID,
		ResourceType: manta.DashboardsResourceType,
		OrgID:        dashboard.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: nil,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add update cell oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", dashboard.ID),
			zap.Stringer("orgID", dashboard.OrgID))
		return err
	}

	return nil
}

func (s *DashboardService) ReplaceDashboardCells(ctx context.Context, dashboardID manta.ID, cells []manta.Cell) error {
	auth, err := authorizer.FromContext(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	err = s.DashboardService.ReplaceDashboardCells(ctx, dashboardID, cells)
	if err != nil {
		return err
	}

	dashboard, err := s.DashboardService.FindDashboardByID(ctx, dashboardID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(dashboard)
	if err != nil {
		return err
	}

	err = s.oplog.AddLogEntry(ctx, manta.OperationLogEntry{
		Type:         manta.Update,
		ResourceID:   dashboard.ID,
		ResourceType: manta.DashboardsResourceType,
		OrgID:        dashboard.OrgID,
		UserID:       auth.GetUserID(),
		ResourceBody: data,
		Time:         now,
	})
	if err != nil {
		s.logger.Error("add update cell oplog failed",
			zap.Error(err),
			zap.Stringer("resourceID", dashboard.ID),
			zap.Stringer("orgID", dashboard.OrgID))
		return err
	}

	return nil
}
