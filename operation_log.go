package manta

import (
	"context"
	"time"
)

const (
	// Create a resource.
	Create ChangeType = "create"
	// Put a resource.
	Put = "put"
	// Update a resource.
	Update = "update"
	// Delete a resource
	Delete = "delete"
)

type ChangeType string

func (c ChangeType) String() string {
	return string(c)
}

// OperationLogEntry to a resource
type OperationLogEntry struct {
	// Type of change
	Type ChangeType `json:"type"`
	// ResourceID of the changed resource
	ResourceID ID `json:"resourceID"`
	// ResourceType that was changed
	ResourceType ResourceType `json:"resourceType"`
	// OrgID of the organization owning the changed resource.
	OrgID ID `json:"orgID"`
	// UserID of the suer who changing the resource.
	UserID ID `json:"userID"`
	// Resourcebody after the change
	ResourceBody []byte `json:"resourceBody"`
	// Time when the resource was changed
	Time time.Time `json:"time"`
}

func (o *OperationLogEntry) Valid() error {
	if !o.ResourceID.Valid() {
		return ErrInvalidResourceID
	}

	if !o.OrgID.Valid() {
		return ErrInvalidOrgID
	}

	if !o.UserID.Valid() {
		return ErrInvalidUserID
	}

	return nil
}

type OperationLogService interface {
	// AddLogEntry add an operation log entry.
	AddLogEntry(ctx context.Context, ent OperationLogEntry) error

	// FindOperationLogsByID return operation logs of a resource.
	FindOperationLogsByID(ctx context.Context, id ID, opts FindOptions) ([]*OperationLogEntry, int, error)

	// FindOperationLogsByUser returns operation logs made by a user.
	FindOperationLogsByUser(ctx context.Context, userID ID, opts FindOptions) ([]*OperationLogEntry, int, error)

	// TODO: add a method to delete log entry!?
}
