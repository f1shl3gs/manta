package pubsub

import "context"

type Handler func(msg []byte)

type Unsubscribe func()

type Pubsub interface {
	Publish(ctx context.Context, topic string, msg interface{}) error

	Subscribe(topic string, handler Handler) (Unsubscribe, error)
}
