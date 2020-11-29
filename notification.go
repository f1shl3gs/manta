package manta

import (
	"context"
	"time"
)

type NotificationFilter struct {
	Start  time.Time
	End    time.Time
	UserID *ID
}

type NotificationService interface {
	FindNotificationByID(ctx context.Context, id ID) (*Notification, error)

	FindNotifications(ctx context.Context, filter NotificationFilter) ([]*Notification, error)

	CreateNotification(ctx context.Context, nf *Notification) error

	// todo: add RetentionPolicy to clean up notifications
}
