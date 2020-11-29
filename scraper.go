package manta

import (
	"context"
)

type ScraperTargetFilter struct {
	OrgID *ID
}

type ScraperTargetUpdate struct {
	Target *string
	Labels *map[string]string
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
