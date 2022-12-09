package broadcast

import (
	"sync"
)

type Broadcaster[T any] struct {
	mtx    sync.RWMutex
	n      uint64
	queues map[uint64]*Queue[T]
}

func New[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		queues: make(map[uint64]*Queue[T]),
	}
}

func (b *Broadcaster[T]) Pub(v T) {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	for _, q := range b.queues {
		select {
		case q.ch <- v:
		default:
		}
	}
}

func (b *Broadcaster[T]) Sub() *Queue[T] {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.n += 1

	q := &Queue[T]{
		id:          b.n,
		ch:          make(chan T, 8),
		broadcaster: b,
	}

	b.queues[b.n] = q

	return q
}

func (b *Broadcaster[T]) Close() {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	for id, q := range b.queues {
		delete(b.queues, id)

		close(q.ch)
	}
}

func (b *Broadcaster[T]) unsub(id uint64) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	delete(b.queues, id)
}

type Queue[T any] struct {
	id          uint64
	ch          chan T
	broadcaster *Broadcaster[T]
}

func (q *Queue[T]) Ch() <-chan T {
	return q.ch
}

func (q *Queue[T]) Close() {
	q.broadcaster.unsub(q.id)
}
