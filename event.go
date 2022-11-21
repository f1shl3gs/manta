package manta

import (
	"context"
	"time"
)

type Acknowledgement struct {
	Username string    `json:"username,omitempty"`
	Desc     string    `json:"desc,omitempty"`
	When     time.Time `json:"when"`
}

type InhibitionStatus struct {
	When         time.Time `json:"when"`
	InhibitionID ID        `json:"inhibitionID,omitempty"`
	Name         string    `json:"name,omitempty"`
	Desc         string    `json:"desc,omitempty"`
}

type EventStatus struct {
	// pending, firing or resolve
	Phase string `json:"phase,omitempty"`
	// acknowledgements
	Acks []Acknowledgement `json:"acks"`
	// inhibitions
	Inhibitions []InhibitionStatus `json:"inhibitions"`
}

type Event struct {
	ID          ID                `json:"id,omitempty"`
	Start       time.Time         `json:"start"`
	End         time.Time         `json:"end"`
	Name        string            `json:"name,omitempty"`
	OrgID       ID                `json:"orgID,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Status      EventStatus       `json:"status"`
}

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
