package manta

import (
	"context"
)

var (
	ErrScraperNotFound = &Error{
		Code: ENotFound,
		Msg:  "scraper not found",
	}
)

type ScraperTargetFilter struct {
	OrgID *ID
}

type ScraperTargetUpdate struct {
	Name    *string
	Desc    *string
	Labels  *map[string]string
	Targets *[]string
}

func (upd *ScraperTargetUpdate) Apply(s *ScrapeTarget) {
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

// ScraperTargetService defines the crud service for ScraperTarget
type ScraperTargetService interface {

	// FindScraperTargetByID returns a single ScraperTarget by ID
	FindScraperTargetByID(ctx context.Context, id ID) (*ScrapeTarget, error)

	// FindScraperTargets returns a list of ScraperTargets that match the filter
	FindScraperTargets(ctx context.Context, filter ScraperTargetFilter) ([]*ScrapeTarget, error)

	// CreateScraperTarget create a ScraperTarget
	CreateScraperTarget(ctx context.Context, target *ScrapeTarget) error

	// UpdateScraperTarget update a single ScraperTarget with chageset
	// returns the new ScraperTarget after update
	UpdateScraperTarget(ctx context.Context, id ID, u ScraperTargetUpdate) (*ScrapeTarget, error)

	// DeleteScraperTarget delete a single ScraperTarget by ID
	DeleteScraperTarget(ctx context.Context, id ID) error
}
