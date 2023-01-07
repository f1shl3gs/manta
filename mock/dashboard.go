package mock

import (
	"context"

	"github.com/f1shl3gs/manta"
)

var _ manta.DashboardService = &DashboardService{}

type DashboardService struct {
	FindDashboardByIDFn func(context.Context, manta.ID) (*manta.Dashboard, error)

	FindDashboardsFn func(context.Context, manta.DashboardFilter) ([]*manta.Dashboard, error)

	CreateDashboardFn func(context.Context, *manta.Dashboard) error

	UpdateDashboardFn func(context.Context, manta.DashboardUpdate) (*manta.Dashboard, error)

	AddDashboardCellFn func(context.Context, manta.ID, *manta.Cell) error

	// RemoveDashboardCell remove a panel by ID
	RemoveDashboardCellFn func(ctx context.Context, dashboardId, cellId manta.ID) error

	// UpdateDashboardCell update the dashboard cell with the provided ids
	UpdateDashboardCellFn func(context.Context, manta.DashboardCellUpdate) (*manta.Cell, error)

	GetDashboardCellFn func(ctx context.Context, dashboardID, cellId manta.ID) (*manta.Cell, error)

	// RemoveDashboard removes dashboard by id
	DeleteDashboardFn func(ctx context.Context, id manta.ID) error

	ReplaceDashboardCellsFn func(ctx context.Context, dashboardId manta.ID, cells []manta.Cell) error
}

func (s *DashboardService) FindDashboardByID(ctx context.Context, id manta.ID) (*manta.Dashboard, error) {
	return s.FindDashboardByIDFn(ctx, id)
}

func (s *DashboardService) FindDashboards(ctx context.Context, filter manta.DashboardFilter) ([]*manta.Dashboard, error) {
	return s.FindDashboardsFn(ctx, filter)
}

func (s *DashboardService) CreateDashboard(ctx context.Context, d *manta.Dashboard) error {
	return s.CreateDashboardFn(ctx, d)
}

func (s *DashboardService) UpdateDashboard(ctx context.Context, upd manta.DashboardUpdate) (*manta.Dashboard, error) {
	return s.UpdateDashboardFn(ctx, upd)
}

func (s *DashboardService) AddDashboardCell(ctx context.Context, id manta.ID, cell *manta.Cell) error {
	return s.AddDashboardCellFn(ctx, id, cell)
}

// RemoveDashboardCell remove a panel by ID
func (s *DashboardService) RemoveDashboardCell(ctx context.Context, dashboardId, cellId manta.ID) error {
	return s.RemoveDashboardCellFn(ctx, dashboardId, cellId)
}

// UpdateDashboardCell update the dashboard cell with the provided ids
func (s *DashboardService) UpdateDashboardCell(ctx context.Context, upd manta.DashboardCellUpdate) (*manta.Cell, error) {
	return s.UpdateDashboardCellFn(ctx, upd)
}

func (s *DashboardService) GetDashboardCell(ctx context.Context, dashboardID, cellId manta.ID) (*manta.Cell, error) {
	return s.GetDashboardCellFn(ctx, dashboardID, cellId)
}

// RemoveDashboard removes dashboard by id
func (s *DashboardService) DeleteDashboard(ctx context.Context, id manta.ID) error {
	return s.DeleteDashboardFn(ctx, id)
}

func (s *DashboardService) ReplaceDashboardCells(ctx context.Context, dashboardId manta.ID, cells []manta.Cell) error {
	return s.ReplaceDashboardCellsFn(ctx, dashboardId, cells)
}
