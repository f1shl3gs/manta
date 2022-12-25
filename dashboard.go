package manta

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Dashboard struct {
	ID      ID        `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Name    string    `json:"name"`
	Desc    string    `json:"desc,omitempty"`
	OrgID   ID        `json:"orgID"`
	Cells   []Cell    `json:"cells,omitempty"`
}

type ViewProperties interface {
	GetType() string
}

type Axis struct {
	Bounds []string `json:"bounds,omitempty"`
	Label  string   `json:"label,omitempty"`
	Prefix string   `json:"prefix,omitempty"`
	Suffix string   `json:"suffix,omitempty"`
	Base   string   `json:"base"`
}

type Axes struct {
	X Axis `json:"x"`
	Y Axis `json:"y"`
}

type Query struct {
	Name   string `json:"name,omitempty"`
	Text   string `json:"text,omitempty"`
	Legend string `json:"legend,omitempty"`
	Hidden bool   `json:"hidden,omitempty"`
}

type DashboardColor struct {
	Id    string `json:"id,omitempty"`
	Type  string `json:"type,omitempty"`
	Hex   string `json:"hex,omitempty"`
	Name  string `json:"name,omitempty"`
	Value int64  `json:"value,omitempty"`
}

type Cell struct {
	ID   ID     `json:"id"`
	Name string `json:"name,omitempty"`
	Desc string `json:"desc,omitempty"`
	X    int32  `json:"x"`
	Y    int32  `json:"y"`
	W    int32  `json:"w"`
	H    int32  `json:"h"`
	MinH int32  `json:"minH,omitempty"`
	MinW int32  `json:"minW,omitempty"`
	MaxW int32  `json:"maxW,omitempty"`

	// Types that are valid to be assigned to ViewProperties:
	//	*Gauge
	//	*XY
	//	*SingleStat
	//	*LinePlusSingleStat
	//	*Markdown
	ViewProperties ViewProperties `json:"viewProperties,omitempty"`
}

type GaugeViewProperties struct {
	Type              string           `json:"type,omitempty"`
	Axes              Axes             `json:"axes"`
	Queries           []Query          `json:"queries,omitempty"`
	Prefix            string           `json:"prefix,omitempty"`
	Suffix            string           `json:"suffix,omitempty"`
	TickPrefix        string           `json:"tickPrefix,omitempty"`
	TickSuffix        string           `json:"tickSuffix,omitempty"`
	Note              string           `json:"note,omitempty"`
	ShowNoteWhenEmpty bool             `json:"showNoteWhenEmpty,omitempty"`
	DecimalPlaces     DecimalPlaces    `json:"decimalPlaces,omitempty"`
	Colors            []DashboardColor `json:"colors,omitempty"`
}

func (g *GaugeViewProperties) GetType() string {
	return g.Type
}

type XYViewProperties struct {
	Type           string           `json:"type,omitempty"`
	Axes           Axes             `json:"axes"`
	Queries        []Query          `json:"queries,omitempty"`
	TimeFormat     string           `json:"timeFormat,omitempty"`
	XColumn        string           `json:"xColumn,omitempty"`
	YColumn        string           `json:"yColumn,omitempty"`
	HoverDimension string           `json:"hoverDimension,omitempty"`
	Position       string           `json:"position,omitempty"`
	Geom           string           `json:"geom,omitempty"`
	Colors         []DashboardColor `json:"colors,omitempty"`
}

func (x *XYViewProperties) GetType() string {
	return x.Type
}

type SingleStatViewProperties struct {
	Type              string           `json:"type,omitempty"`
	Note              string           `json:"note,omitempty"`
	Queries           []Query          `json:"queries"`
	Prefix            string           `json:"prefix,omitempty"`
	Suffix            string           `json:"suffix,omitempty"`
	TickPrefix        string           `json:"tickPrefix,omitempty"`
	TickSuffix        string           `json:"tickSuffix,omitempty"`
	ShowNoteWhenEmpty bool             `json:"showNoteWhenEmpty,omitempty"`
	Colors            []DashboardColor `json:"colors"`
	DecimalPlaces     DecimalPlaces    `json:"decimalPlaces,omitempty"`
}

func (s *SingleStatViewProperties) GetType() string {
	return s.Type
}

type DecimalPlaces struct {
	IsEnforced bool  `json:"isEnforced,omitempty"`
	Digits     int32 `json:"digits,omitempty"`
}

type LinePlusSingleStatViewProperties struct {
	Type              string           `json:"type,omitempty"`
	Note              string           `json:"note,omitempty"`
	Queries           []Query          `json:"queries"`
	Prefix            string           `json:"prefix,omitempty"`
	Suffix            string           `json:"suffix,omitempty"`
	TickPrefix        string           `json:"tickPrefix,omitempty"`
	TickSuffix        string           `json:"tickSuffix,omitempty"`
	ShowNoteWhenEmpty bool             `json:"showNoteWhenEmpty,omitempty"`
	DecimalPlaces     DecimalPlaces    `json:"decimalPlaces"`
	Axes              Axes             `json:"axes"`
	Colors            []DashboardColor `json:"colors"`
}

func (s *LinePlusSingleStatViewProperties) GetType() string {
	return s.Type
}

type MarkdownViewProperties struct {
	Type    string `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
}

func (s *MarkdownViewProperties) GetType() string {
	return s.Type
}

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
	ViewProperties ViewProperties
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

func unmarshalCellPropertiesJSON(b []byte) (ViewProperties, error) {
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
