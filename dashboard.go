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

func (m *Panel) UnmarshalJSON(b []byte) error {
	var p struct {
		Name        string          `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
		Description string          `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
		W           uint32          `protobuf:"varint,3,opt,name=w,proto3" json:"w,omitempty"`
		H           uint32          `protobuf:"varint,4,opt,name=h,proto3" json:"h,omitempty"`
		X           uint32          `protobuf:"varint,5,opt,name=x,proto3" json:"x,omitempty"`
		Y           uint32          `protobuf:"varint,6,opt,name=y,proto3" json:"y,omitempty"`
		Queries     []Query         `protobuf:"bytes,7,rep,name=queries,proto3" json:"queries"`
		Properties  json.RawMessage `json:"properties"`
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	// set values
	m.Name = p.Name
	m.Description = p.Description
	m.W = p.W
	m.H = p.H
	m.X = p.X
	m.Y = p.Y
	m.Queries = p.Queries

	props, err := unmarshalPanelPropertiesJSON(p.Properties)
	if err != nil {
		return err
	}

	m.Properties = props
	return nil
}

func unmarshalPanelPropertiesJSON(b []byte) (isPanel_Properties, error) {
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
