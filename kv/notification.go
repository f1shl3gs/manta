package kv

import (
	"context"

	"github.com/f1shl3gs/manta"
)

var (
	notificationBucket          = []byte("notification")
	notificationUserIndexBucket = []byte("notificationuserindex")
)

var _ manta.NotificationService = (*Service)(nil)

func (s *Service) FindNotificationByID(ctx context.Context, id manta.ID) (*manta.Notification, error) {

	panic("implement me")
}

func (s *Service) FindNotifications(ctx context.Context, filter manta.NotificationFilter) ([]*manta.Notification, error) {
	panic("implement me")
}

func (s *Service) CreateNotification(ctx context.Context, nf *manta.Notification) error {
	panic("implement me")
}
