package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
)

var (
	eventBucket             = []byte("event")
	eventOrgIndexBucket     = []byte("eventorgindex")
	eventOrgNameIndexBucket = []byte("eventorgnameindex")
)

func (s *Service) FindEventByID(ctx context.Context, id manta.ID) (*manta.Event, error) {
	var (
		ev  *manta.Event
		err error
	)

	err = s.kv.View(ctx, func(tx Tx) error {
		ev, err = s.findEventByID(ctx, tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	return ev, nil
}

func (s *Service) findEventByID(ctx context.Context, tx Tx, id manta.ID) (*manta.Event, error) {
	pk, err := id.Encode()
	if err != nil {
		return nil, err
	}

	b, err := tx.Bucket(eventBucket)
	if err != nil {
		return nil, err
	}

	data, err := b.Get(pk)
	if err != nil {
		return nil, err
	}

	ev := &manta.Event{}
	if err = ev.Unmarshal(data); err != nil {
		return nil, err
	}

	return ev, nil
}

func (s *Service) FindEvents(ctx context.Context, filter manta.EventFilter, opt ...manta.FindOptions) ([]*manta.Event, int, error) {
	panic("implement me")
}

func (s *Service) CreateEvent(ctx context.Context, ev *manta.Event) error {
	panic("implement me")
}

func (s *Service) putEvent(ctx context.Context, tx Tx, ev *manta.Event) error {
	return nil
}

func (s *Service) UpdateEvent(ctx context.Context, id manta.ID, u manta.UpdateEvent) (*manta.Event, error) {
	panic("implement me")
}

func (s *Service) updateEvent(ctx context.Context, tx Tx, id manta.ID, u manta.UpdateEvent) (*manta.Event, error) {
	panic("implement me")
}

func (s *Service) DeleteEvent(ctx context.Context, id manta.ID) error {
	panic("implement me")
}

func (s *Service) deleteEvent(ctx context.Context, tx Tx, id manta.ID) error {
	return nil
}
