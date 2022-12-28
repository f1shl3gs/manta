package manta

import (
	"context"
	"encoding/json"
)

// NotificationEndpoint is the configuration describing
// how to call a 3rd party service.
type NotificationEndpoint interface {
	json.Marshaler
	CRUDGetter
	CRUDSetter

	Type() string
	Valid() error

	GetID() ID
	SetID(id ID)

	GetOrgID() ID
	SetOrgID(orgID ID)

	GetName() string
	SetName(name string)

	GetDesc() string
	SetDesc(desc string)

	// BackfillSecretKeys fill back fill the secret field key during the unmarshalling
	// if value of that secret field is not nil.
	BackfillSecretKeys()
	// SecretFields return available secret fields.
	SecretFields() []SecretField
}

type NotificationEndpointFilter struct {
	OrgID ID
}

type NotificationEndpointUpdate struct {
	Name *string `json:"name"`
	Desc *string `json:"desc"`
}

func (upd *NotificationEndpointUpdate) Apply(ne NotificationEndpoint) {
	if upd.Name != nil {
		ne.SetName(*upd.Name)
	}

	if upd.Desc != nil {
		ne.SetDesc(*upd.Desc)
	}
}

type NotificationEndpointService interface {
	// FindNotificationByID returns a single notification endpoint by ID
	FindNotificationEndpointByID(ctx context.Context, id ID) (NotificationEndpoint, error)

	// FindNotificationEndpoints returns a list of notication endpoints that match filter.
	FindNotificationEndpoints(ctx context.Context, filter NotificationEndpointFilter) ([]NotificationEndpoint, error)

	// CreateNotificationEndpoint creates a new notification endpoint and sets b.ID with the new identifier
	CreateNotificationEndpoint(ctx context.Context, ne NotificationEndpoint) error

	// UpdateNotificationEndpoint updates a single notification endpoint.
	// Returns the new notification endpoint after update.
	UpdateNotificationEndpoint(ctx context.Context, id ID, ne NotificationEndpoint) (NotificationEndpoint, error)

	// PatchNotificationENdpoint patch a single notification endpoint.
	// Returns the new notification endpoint after patch.
	PatchNotificationEndpoint(ctx context.Context, id ID, upd NotificationEndpointUpdate) (NotificationEndpoint, error)

	// DeleteNotificationEndpoint remove a notification endpoint by ID, return it's secret fields, orgID for further deletion
	DeleteNotificationEndpoint(ctx context.Context, id ID) ([]SecretField, ID, error)
}
