package manta

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/influxdata/cron"
)

type Threshold struct {
	Type  string  `json:"type"`
	Value float64 `json:"value,omitempty"`
	// for inside and outside only
	Min float64 `json:"min,omitempty"`
	Max float64 `json:"max,omitempty"`
}

type Condition struct {
	Status    string        `json:"status"`
	Pending   time.Duration `json:"pending"`
	Threshold Threshold     `json:"threshold"`
}

func (m *Condition) Validate() error {
	if m.Threshold.Type != "inside" && m.Threshold.Type != "outside" {
		return nil
	}

	switch m.Threshold.Type {
	case "inside", "outside":
		if m.Threshold.Max <= m.Threshold.Min {
			return invalidField("max", errors.New("condition.max must be larger than min"))
		}

	default:
	}

	return nil
}

type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Check struct {
	ID         ID          `json:"id"`
	Created    time.Time   `json:"created"`
	Updated    time.Time   `json:"updated"`
	Name       string      `json:"name,omitempty"`
	Desc       string      `json:"desc,omitempty"`
	OrgID      ID          `json:"orgID,omitempty"`
	Expr       string      `json:"expr,omitempty"`
	Status     string      `json:"status,omitempty"`
	Cron       string      `json:"cron,omitempty"`
	Conditions []Condition `json:"conditions"`
	TaskID     ID          `json:"taskId"`
	Labels     []Label     `json:"labels,omitempty"`
}

func (c *Check) GetID() ID {
	return c.ID
}

func (c *Check) GetOrgID() ID {
	return c.OrgID
}

func (c *Check) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *Check) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (m *Check) Validate() error {
	if m.Name == "" {
		return invalidField("name", ErrFieldMustBeSet)
	}

	if m.Desc == "" {
		return invalidField("desc", ErrFieldMustBeSet)
	}

	if m.Expr == "" {
		return invalidField("expr", ErrFieldMustBeSet)
	}

	if m.Status == "" {
		return invalidField("status", ErrFieldMustBeSet)
	}

	if err := validateStatus(m.Status); err != nil {
		return invalidField("status", err)
	}

	if _, err := cron.ParseUTC(m.Cron); err != nil {
		return invalidField("cron", err)
	}

	if len(m.Conditions) == 0 {
		return invalidField("conditions", ErrFieldMustBeSet)
	}

	for i := 0; i < len(m.Conditions); i++ {
		if err := m.Conditions[i].Validate(); err != nil {
			return err
		}
	}

	return nil
}

type CheckUpdate struct {
	Name   *string `json:"name"`
	Desc   *string `json:"desc"`
	Status *string `json:"status"`
}

func (upd *CheckUpdate) Apply(check *Check) {
	if upd.Name != nil {
		check.Name = *upd.Name
	}

	if upd.Desc != nil {
		check.Desc = *upd.Desc
	}

	if upd.Status != nil {
		check.Status = *upd.Status
	}
}

type CheckFilter struct {
	OrgID *ID
}

type CheckService interface {
	// FindCheckByID returns a check by id
	FindCheckByID(ctx context.Context, id ID) (*Check, error)

	// FindChecks returns a list of checks that match the filter and total count of matching checks
	// Additional options provide pagination & sorting.
	FindChecks(ctx context.Context, filter CheckFilter, opt ...FindOptions) ([]*Check, int, error)

	// CreateCheck creates a new and set its id with new identifier
	CreateCheck(ctx context.Context, c *Check) error

	// UpdateCheck updates the whole check
	// Returns the new check after update
	UpdateCheck(ctx context.Context, id ID, c *Check) (*Check, error)

	// PatchCheck updates a single check with changeset
	// Returns the new check after patch
	PatchCheck(ctx context.Context, id ID, u CheckUpdate) (*Check, error)

	// DeleteCheck delete a single check by ID
	DeleteCheck(ctx context.Context, id ID) error
}

// ThresholdType is the Condition's Threshold
type ThresholdType string

const (
	NoDate    = "nodata"
	GreatThan = "gt"
	Equal     = "eq"
	NotEqual  = "ne"
	LessThan  = "lt"
	Inside    = "inside"
	Outside   = "outside"
)