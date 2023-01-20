package manta

import (
	"context"
	"time"
)

type ScrapeTarget struct {
	ID      ID                `json:"id,omitempty"`
	Created time.Time         `json:"created"`
	Updated time.Time         `json:"updated"`
	OrgID   ID                `json:"orgID,omitempty"`
	Name    string            `json:"name,omitempty"`
	Desc    string            `json:"desc,omitempty"`
	Targets []string          `json:"targets,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func (s *ScrapeTarget) GetID() ID {
	return s.ID
}

func (s *ScrapeTarget) GetOrgID() ID {
	return s.OrgID
}

type ScrapeTargetFilter struct {
	OrgID *ID
}

type ScrapeTargetUpdate struct {
	OrgID ID

	Name    *string
	Desc    *string
	Labels  *map[string]string
	Targets *[]string
}

func (upd *ScrapeTargetUpdate) Apply(s *ScrapeTarget) {
	if upd.Name != nil {
		s.Name = *upd.Name
	}

	if upd.Desc != nil {
		s.Desc = *upd.Desc
	}

	if upd.Labels != nil {
		s.Labels = *upd.Labels
	}

	if upd.Targets != nil {
		s.Targets = *upd.Targets
	}
}

func (m *ScrapeTarget) Validate() error {
	if m.Name == "" {
		return &Error{Code: EInvalid, Msg: "Name is required"}
	}

	if !m.OrgID.Valid() {
		return ErrInvalidOrgID
	}

	return nil
}

// ScrapeTargetService defines the crud service for ScraperTarget
type ScrapeTargetService interface {
	// FindScrapeTargetByID returns a single ScraperTarget by ID
	FindScrapeTargetByID(ctx context.Context, id ID) (*ScrapeTarget, error)

	// FindScrapeTargets returns a list of ScraperTargets that match the filter
	FindScrapeTargets(ctx context.Context, filter ScrapeTargetFilter) ([]*ScrapeTarget, error)

	// CreateScrapeTarget create a ScraperTarget
	CreateScrapeTarget(ctx context.Context, target *ScrapeTarget) error

	// UpdateScrapeTarget update a single ScraperTarget with chageset
	// returns the new ScraperTarget after update
	UpdateScrapeTarget(ctx context.Context, id ID, upd ScrapeTargetUpdate) (*ScrapeTarget, error)

	// DeleteScrapeTarget delete a single ScraperTarget by ID
	DeleteScrapeTarget(ctx context.Context, id ID) error
}
