package manta

import "context"

type DashboardFilter struct {
	OrganizationID *ID
}

type DashboardUpdate struct {
	Name *string
	Desc *string
}

type DashboardPanelUpdate struct {
	W, H, X, Y *uint32
}

type DashboardService interface {
	FindDashboardByID(ctx context.Context, id ID) (*Dashboard, error)

	FindDashboards(ctx context.Context, filter DashboardFilter) ([]*Dashboard, error)

	CreateDashboard(ctx context.Context, d *Dashboard) error

	UpdateDashboard(ctx context.Context, udp DashboardUpdate) (*Dashboard, error)

	AddDashboardPanel(ctx context.Context, p Panel) error

	// RemoveDashboardPanel remove a panel by ID
	RemoveDashboardPanel(ctx context.Context, did, pid ID) error

	// UpdateDashboardPanel update the dashboard panel with the provided ids
	UpdateDashboardPanel(ctx context.Context, did, pid ID, udp DashboardPanelUpdate) (Panel, error)

	// RemoveDashboard removes dashboard by id
	RemoveDashboard(ctx context.Context, id ID) error
}
