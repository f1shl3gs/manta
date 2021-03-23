package manta

import (
	"context"
	"encoding/json"

	"github.com/f1shl3gs/manta/pkg/duration"
)

const (
	OK   = "ok"
	Info = "info"
	Warn = "warn"
	High = "high"
	Crit = "crit"
)

var SeverityValue = map[string]int{
	OK:   0,
	Info: 1,
	Warn: 2,
	High: 3,
	Crit: 4,
}

var Severities = []string{
	OK,
	Info,
	Warn,
	High,
	Crit,
}

type CheckFilter struct {
	OrgID *ID
}

type CheckUpdate struct {
	Name   *string `json:"name"`
	Desc   *string `json:"desc"`
	Status *string `json:"status"`
}

func (upd *CheckUpdate) Validate() error {
	if upd.Name != nil {
		if *upd.Name == "" {
			return &Error{Code: EInvalid, Op: "validate check's Name", Msg: "Name cannot be empty"}
		}
	}

	if upd.Status != nil {
		if *upd.Status != "active" && *upd.Status != "inactive" {
			return &Error{Code: EInvalid, Op: "validate check's Status", Msg: "status is not active nor inactive"}
		}
	}

	return nil
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

var (
	thresholdTypes = []string{
		NoDate,
		GreatThan,
		Equal,
		NotEqual,
		LessThan,
		Inside,
		Outside,
	}
)

func (m *Condition) UnmarshalJSON(b []byte) error {
	var tm struct {
		Status    string    `json:"status"`
		Pending   string    `json:"pending"`
		Threshold Threshold `json:"threshold"`
	}

	if err := json.Unmarshal(b, &tm); err != nil {
		return err
	}

	if tm.Pending != "" {
		d, err := duration.Parse(tm.Pending)
		if err != nil {
			return err
		}

		m.Pending = d
	}

	m.Status = tm.Status
	m.Threshold = tm.Threshold

	return nil
}

func (m *Condition) MarshalJSON() ([]byte, error) {
	var tm struct {
		Status    string    `json:"status"`
		Pending   string    `json:"pending"`
		Threshold Threshold `json:"threshold"`
	}

	tm.Status = m.Status
	tm.Pending = m.Pending.String()
	tm.Threshold = m.Threshold

	return json.Marshal(tm)
}
