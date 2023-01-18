package raftstore

import (
	"fmt"
	"sync"
)

const (
	// To avoid lock contention we use an array of wait struct (rw mutex & map)
	// for the id argument, we apply mod operation and uses its remainder to
	// index into the array and find the corresponding element.
	defaultListElementLength = 64
)

// wait is an interface that provides the ability to wait and trigger events that
// are associated with IDs.
type wait struct {
	shards []waitShard
}

type waitShard struct {
	l sync.RWMutex
	m map[uint64]chan interface{}
}

// newWait creates a Wait.
func newWait() *wait {
	res := wait{
		shards: make([]waitShard, defaultListElementLength),
	}
	for i := 0; i < len(res.shards); i++ {
		res.shards[i].m = make(map[uint64]chan interface{})
	}
	return &res
}

// Register waits returns a chan that waits on the given ID.
// The chan will be triggered when Trigger is called with
// the same ID.
func (w *wait) Register(id uint64) <-chan interface{} {
	idx := id % defaultListElementLength
	newCh := make(chan interface{}, 1)

	w.shards[idx].l.Lock()
	_, ok := w.shards[idx].m[id]
	if !ok {
		w.shards[idx].m[id] = newCh
	}
	w.shards[idx].l.Unlock()

	if ok {
		panic(fmt.Sprintf("duplicated wait id %x", id))
	}

	return newCh
}

// Trigger triggers the waiting chans with the given ID.
func (w *wait) Trigger(id uint64, x interface{}) {
	idx := id % defaultListElementLength

	w.shards[idx].l.Lock()
	ch := w.shards[idx].m[id]
	delete(w.shards[idx].m, id)
	w.shards[idx].l.Unlock()

	if ch != nil {
		ch <- x
		close(ch)
	}
}

func (w *wait) IsRegistered(id uint64) bool {
	idx := id % defaultListElementLength

	w.shards[idx].l.RLock()
	_, ok := w.shards[idx].m[id]
	w.shards[idx].l.RUnlock()

	return ok
}

type waitTime struct {
	mtx                 sync.Mutex
	lastTriggerDeadline uint64
	m                   map[uint64]chan struct{}
	closedCh            chan struct{}
}

func newWaitTime() *waitTime {
	wt := &waitTime{
		m:        make(map[uint64]chan struct{}),
		closedCh: make(chan struct{}),
	}

	close(wt.closedCh)

	return wt
}

func (w *waitTime) Wait(deadline uint64) <-chan struct{} {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.lastTriggerDeadline >= deadline {
		return w.closedCh
	}
	ch := w.m[deadline]
	if ch == nil {
		ch = make(chan struct{})
		w.m[deadline] = ch
	}
	return ch
}

func (w *waitTime) Trigger(deadline uint64) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	w.lastTriggerDeadline = deadline
	for t, ch := range w.m {
		if t <= deadline {
			delete(w.m, t)
			close(ch)
		}
	}
}
