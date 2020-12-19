package manta

import (
	"context"
	"encoding/json"
	"errors"
)

type DashboardFilter struct {
	OrganizationID *ID
}

type DashboardUpdate struct {
	Name *string
	Desc *string
}

func (udp DashboardUpdate) Apply(dash *Dashboard) {
	if udp.Name != nil {
		dash.Name = *udp.Name
	}

	if udp.Desc != nil {
		dash.Desc = *udp.Name
	}
}

type DashboardCellUpdate struct {
	W, H, X, Y *int32
}

func (udp DashboardCellUpdate) Apply(cell *Cell) {
	if udp.W != nil {
		cell.W = *udp.W
	}

	if udp.H != nil {
		cell.H = *udp.H
	}

	if udp.X != nil {
		cell.X = *udp.X
	}

	if udp.Y != nil {
		cell.Y = *udp.Y
	}
}

type DashboardService interface {
	FindDashboardByID(ctx context.Context, id ID) (*Dashboard, error)

	FindDashboards(ctx context.Context, filter DashboardFilter) ([]*Dashboard, error)

	CreateDashboard(ctx context.Context, d *Dashboard) error

	UpdateDashboard(ctx context.Context, id ID, udp DashboardUpdate) (*Dashboard, error)

	AddDashboardCell(ctx context.Context, id ID, cell *Cell) error

	// RemoveDashboardCell remove a panel by ID
	RemoveDashboardCell(ctx context.Context, did, pid ID) error

	// UpdateDashboardCell update the dashboard cell with the provided ids
	UpdateDashboardCell(ctx context.Context, did, pid ID, udp DashboardCellUpdate) (*Cell, error)

	// RemoveDashboard removes dashboard by id
	DeleteDashboard(ctx context.Context, id ID) error
}

func (m *Cell) UnmarshalJSON(b []byte) error {
	var c struct {
		Name        string          `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
		Description string          `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
		W           int32           `protobuf:"varint,3,opt,name=w,proto3" json:"w,omitempty"`
		H           int32           `protobuf:"varint,4,opt,name=h,proto3" json:"h,omitempty"`
		X           int32           `protobuf:"varint,5,opt,name=x,proto3" json:"x,omitempty"`
		Y           int32           `protobuf:"varint,6,opt,name=y,proto3" json:"y,omitempty"`
		Queries     []Query         `protobuf:"bytes,7,rep,name=queries,proto3" json:"queries"`
		Properties  json.RawMessage `json:"properties"`
	}

	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}

	// set values
	m.Name = c.Name
	m.Description = c.Description
	m.W = c.W
	m.H = c.H
	m.X = c.X
	m.Y = c.Y

	props, err := unmarshalCellPropertiesJSON(c.Properties)
	if err != nil {
		return err
	}

	m.Properties = props
	return nil
}

func unmarshalCellPropertiesJSON(b []byte) (isCell_Properties, error) {
	var t struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	switch t.Type {
	case "xy":
		var xy XYView
		if err := json.Unmarshal(b, &xy); err != nil {
			return nil, err
		}

		return &xy, nil
	default:
		return nil, errors.New("unknown type")
	}
}
