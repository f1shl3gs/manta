package manta

import (
	"context"
	"time"
)

type EventFilter struct {
	Name   *string
	Org    *ID
	Status *string

	Start *time.Time
	End   *time.Time
}

type UpdateEvent struct {
	End         *time.Time
	Labels      map[string]string
	Annotations map[string]string
}

type EventService interface {
	// FindEventByID find a single Event by id
	FindEventByID(ctx context.Context, id ID) (*Event, error)

	// FindEvents return a list of events that match filter and the total count of matching events
	FindEvents(ctx context.Context, filter EventFilter, opt ...FindOptions) ([]*Event, int, error)

	// CreateEvent a single event and sets it's id with the new identifier
	CreateEvent(ctx context.Context, ev *Event) error

	// UpdateEvent update a single event with changeset
	// returns the new event state after update
	UpdateEvent(ctx context.Context, id ID, u UpdateEvent) (*Event, error)

	// DeleteEvent remove an event by ID
	DeleteEvent(ctx context.Context, id ID) error

	// todo: Retention Policy
}
