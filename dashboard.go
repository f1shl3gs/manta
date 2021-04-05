package manta

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

type DashboardFilter struct {
	OrganizationID *ID
}

type DashboardUpdate struct {
	Name *string `json:"name,omitempty"`
	Desc *string `json:"desc,omitempty"`
}

func (udp DashboardUpdate) Apply(dash *Dashboard) {
	if udp.Name != nil {
		dash.Name = *udp.Name
	}

	if udp.Desc != nil {
		dash.Desc = *udp.Desc
	}
}

type DashboardCellUpdate struct {
	Name           *string
	Desc           *string
	W, H, X, Y     *int32
	ViewProperties isCell_ViewProperties
}

func (udp DashboardCellUpdate) Apply(cell *Cell) {
	if udp.Name != nil {
		cell.Name = *udp.Name
	}

	if udp.Desc != nil {
		cell.Desc = *udp.Desc
	}

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

	if udp.ViewProperties != nil {
		cell.ViewProperties = udp.ViewProperties
	}
}

func (udp *DashboardCellUpdate) UnmarshalJSON(bytes []byte) error {
	var a struct {
		Name           *string
		Desc           *string
		W, H, X, Y     *int32
		ViewProperties json.RawMessage `json:"viewProperties"`
	}
	{
	}

	err := json.Unmarshal(bytes, &a)
	if err != nil {
		return err
	}

	vp, err := unmarshalCellPropertiesJSON(a.ViewProperties)
	if err != nil {
		return err
	}

	udp.Name = a.Name
	udp.Desc = a.Desc
	udp.W = a.W
	udp.H = a.H
	udp.X = a.X
	udp.Y = a.Y
	udp.ViewProperties = vp

	return nil
}

type ViewUpdate struct {
}

func (udp ViewUpdate) Apply(view *View) {
	// todo: implement it
}

type View struct {
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

	GetDashboardCell(ctx context.Context, did, cid ID) (*Cell, error)

	// GetDashboardCellView(ctx context.Context, did, cid ID) (*View, error)
	//
	// UpdateDashboardCellView(ctx context.Context, did, cid ID, udp ViewUpdate) (*View, error)

	// RemoveDashboard removes dashboard by id
	DeleteDashboard(ctx context.Context, id ID) error

	ReplaceDashboardCells(ctx context.Context, did ID, cells []Cell) error
}

func (m *Cell) Validate() error {
	if m.ID == 0 {
		return ErrInvalidID
	}

	if m.W == 0 {
		return errors.New("invalid width")
	}

	if m.H == 0 {
		return errors.New("invalid height")
	}

	return nil
}

func (m *Cell) UnmarshalJSON(b []byte) error {
	var c struct {
		ID             ID
		Name           string
		Desc           string
		W              int32
		H              int32
		X              int32
		Y              int32
		ViewProperties json.RawMessage
	}

	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}

	// set values
	m.ID = c.ID
	m.Name = c.Name
	m.Desc = c.Desc
	m.W = c.W
	m.H = c.H
	m.X = c.X
	m.Y = c.Y

	if c.ViewProperties == nil {
		return nil
	}

	props, err := unmarshalCellPropertiesJSON(c.ViewProperties)
	if err != nil {
		return err
	}

	if validator, ok := props.(Validator); ok {
		err = validator.Validate()
		if err != nil {
			return errors.Wrap(err, "invalid ViewProperties")
		}
	}

	m.ViewProperties = props
	return nil
}

type Validator interface {
	Validate() error
}

func unmarshalCellPropertiesJSON(b []byte) (isCell_ViewProperties, error) {
	var t struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	switch t.Type {
	case "xy":
		var xy XYViewProperties
		if err := json.Unmarshal(b, &xy); err != nil {
			return nil, err
		}

		return &xy, nil
	case "gauge":
		var gauge GaugeViewProperties
		if err := json.Unmarshal(b, &gauge); err != nil {
			return nil, err
		}

		return &gauge, nil

	case "single-stat":
		var singleStat SingleStatViewProperties
		if err := json.Unmarshal(b, &singleStat); err != nil {
			return nil, err
		}

		return &singleStat, nil
	case "line-plus-single-stat":
		var lpss LinePlusSingleStatViewProperties
		if err := json.Unmarshal(b, &lpss); err != nil {
			return nil, err
		}

		return &lpss, nil

	default:
		return nil, errors.New("unknown viewProperties type")
	}
}

func (m *XYViewProperties) Validate() error {
	if len(m.Queries) == 0 {
		return errors.New("queries of XYViewProperties is required")
	}

	return nil
}

func (m *GaugeViewProperties) Validate() error {
	if len(m.Queries) == 0 {
		return errors.New("queries is required")
	}

	return nil
}

func (m *SingleStatViewProperties) Validate() error {
	if len(m.Queries) == 0 {
		return errors.New("queries is required")
	}

	return nil
}

func (m *LinePlusSingleStatViewProperties) Validate() error {
	if len(m.Queries) == 0 {
		return errors.New("queries is required")
	}

	return nil
}
