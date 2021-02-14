package pubsub

import "context"

type Handler func(msg interface{})

type Unsubscribe func()

type Pubsub interface {
	Publish(ctx context.Context, topic string, msg interface{}) error

	Subscribe(topic string, handler Handler) (Unsubscribe, error)

	Close()
}
