package manta

import (
	"context"
)

type TemplateFilter struct {
	Name *string
}

type TemplateUpdate struct {
	Name        *string
	Description *string
}

type TemplateService interface {
	// FindTemplateByID returns a single template by id
	FindTemplateByID(ctx context.Context, id ID) (*Template, error)

	// FindTemplateByName returns the template by name
	FindTemplateByName(ctx context.Context, name string) (*Template, error)

	// FindTemplates returns a list of templates that match the filter and the total count of matching servers
	// additional options provide pagination & sorting
	FindTemplates(ctx context.Context, filter TemplateFilter, opt ...FindOptions) ([]*Template, int, error)

	// CreateTemplate create a single template and set it's id with identifier
	CreateTemplate(ctx context.Context, template *Template) error

	// UpdateTemplate update a single template with changeset
	// returns the new template after update
	UpdateTemplate(ctx context.Context, id ID, u TemplateUpdate) (*Template, error)

	// DeleteTemplate delete a single template by id
	DeleteTemplate(ctx context.Context, id ID) error
}
