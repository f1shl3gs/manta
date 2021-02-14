package pubsub

import (
	"context"
	"sync"
	"sync/atomic"
)

type Memory struct {
	idGen    uint64
	mtx      sync.RWMutex
	channels map[string]map[uint64]chan interface{}
}

func (m *Memory) Publish(ctx context.Context, topic string, msg interface{}) error {
	panic("implement me")
}

func (m *Memory) Subscribe(topic string, handler Handler) (Unsubscribe, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan interface{}, 10)
	sid := atomic.AddUint64(&m.idGen, 1)

	ss := m.channels[topic]
	if ss == nil {
		ss = make(map[uint64]chan interface{})
		m.channels[topic] = ss
	}

	ss[sid] = ch

	go func() {
		for {
			select {
			case v := <-ch:
				handler(v)
			case <-ctx.Done():
				return
			}
		}
	}()

	return Unsubscribe(cancel), nil
}

func (m *Memory) Close() {

}
