package authorizer

import (
	"context"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/errors"
)

var _ manta.DashboardService = &DashboardService{}

type DashboardService struct {
	dashboardService manta.DashboardService
}

func NewDashboardService(ds manta.DashboardService) manta.DashboardService {
	return &DashboardService{
		dashboardService: ds,
	}
}

func (s *DashboardService) FindDashboardByID(ctx context.Context, id manta.ID) (*manta.Dashboard, error) {
	d, err := s.dashboardService.FindDashboardByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeRead(ctx, manta.DashboardsResourceType, id, d.OrgID); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *DashboardService) FindDashboards(ctx context.Context, filter manta.DashboardFilter) ([]*manta.Dashboard, error) {
	dashboards, err := s.dashboardService.FindDashboards(ctx, filter)
	if err != nil {
		return nil, err
	}

	filtered := dashboards[:0]
	for _, d := range dashboards {
		_, _, err := authorizeRead(ctx, manta.DashboardsResourceType, d.ID, d.OrgID)
		if err != nil && errors.ErrorCode(err) != errors.EUnauthorized {
			return nil, err
		}

		if errors.ErrorCode(err) == errors.EUnauthorized {
			continue
		}

		filtered = append(filtered, d)
	}

	return filtered, nil
}

func (s *DashboardService) CreateDashboard(ctx context.Context, d *manta.Dashboard) error {
	if _, _, err := authorizeCreate(ctx, manta.DashboardsResourceType, d.OrgID); err != nil {
		return err
	}

	return s.dashboardService.CreateDashboard(ctx, d)
}

func (s *DashboardService) UpdateDashboard(ctx context.Context, upd manta.DashboardUpdate) (*manta.Dashboard, error) {
	if _, _, err := authorizeWrite(ctx, manta.DashboardsResourceType, upd.ID, upd.OrgID); err != nil {
		return nil, err
	}

	return s.dashboardService.UpdateDashboard(ctx, upd)
}

func (s *DashboardService) AddDashboardCell(ctx context.Context, id manta.ID, cell *manta.Cell) error {
	d, err := s.FindDashboardByID(ctx, id)
	if err != nil {
		return err
	}

	if _, _, err := authorizeWrite(ctx, manta.DashboardsResourceType, id, d.OrgID); err != nil {
		return err
	}

	return s.dashboardService.AddDashboardCell(ctx, id, cell)
}

// RemoveDashboardCell remove a panel by ID
func (s *DashboardService) RemoveDashboardCell(ctx context.Context, dashboardID, cellID manta.ID) error {
	d, err := s.FindDashboardByID(ctx, dashboardID)
	if err != nil {
		return err
	}

	if _, _, err := authorizeWrite(ctx, manta.DashboardsResourceType, dashboardID, d.OrgID); err != nil {
		return err
	}

	return s.dashboardService.RemoveDashboardCell(ctx, dashboardID, cellID)
}

// UpdateDashboardCell update the dashboard cell with the provided ids
func (s *DashboardService) UpdateDashboardCell(ctx context.Context, upd manta.DashboardCellUpdate) (*manta.Cell, error) {
	d, err := s.FindDashboardByID(ctx, upd.DashboardID)
	if err != nil {
		return nil, err
	}

	if _, _, err := authorizeWrite(ctx, manta.DashboardsResourceType, upd.DashboardID, d.OrgID); err != nil {
		return nil, err
	}

	return s.dashboardService.UpdateDashboardCell(ctx, upd)
}

func (s *DashboardService) GetDashboardCell(ctx context.Context, dashboardID, cellID manta.ID) (*manta.Cell, error) {
	_, err := s.dashboardService.GetDashboardCell(ctx, dashboardID, cellID)
	if err != nil {
		return nil, err
	}

	panic("not implement")
}

// DeleteDashboard removes dashboard by id
func (s *DashboardService) DeleteDashboard(ctx context.Context, id manta.ID) error {
	d, err := s.FindDashboardByID(ctx, id)
	if err != nil {
		return err
	}

	if _, _, err = authorizeWrite(ctx, manta.DashboardsResourceType, d.ID, d.OrgID); err != nil {
		return err
	}

	return s.dashboardService.DeleteDashboard(ctx, id)
}

func (s *DashboardService) ReplaceDashboardCells(ctx context.Context, did manta.ID, cells []manta.Cell) error {
	d, err := s.FindDashboardByID(ctx, did)
	if err != nil {
		return err
	}

	if _, _, err = authorizeWrite(ctx, manta.DashboardsResourceType, d.ID, d.OrgID); err != nil {
		return err
	}

	return s.dashboardService.ReplaceDashboardCells(ctx, did, cells)
}
