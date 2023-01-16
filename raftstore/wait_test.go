package raftstore

import (
	"fmt"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	const eid = 1
	wt := newWait()
	ch := wt.Register(eid)
	wt.Trigger(eid, "foo")
	v := <-ch
	if g, w := fmt.Sprintf("%v (%T)", v, v), "foo (string)"; g != w {
		t.Errorf("<-ch = %v, want %v", g, w)
	}

	if g := <-ch; g != nil {
		t.Errorf("unexpected non-nil value: %v (%T)", g, g)
	}
}

func TestRegisterDupPanic(t *testing.T) {
	const eid = 1
	wt := newWait()
	ch1 := wt.Register(eid)

	panicC := make(chan struct{}, 1)

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicC <- struct{}{}
			}
		}()
		wt.Register(eid)
	}()

	select {
	case <-panicC:
	case <-time.After(1 * time.Second):
		t.Errorf("failed to receive panic")
	}

	wt.Trigger(eid, "foo")
	<-ch1
}

func TestTriggerDupSuppression(t *testing.T) {
	const eid = 1
	wt := newWait()
	ch := wt.Register(eid)
	wt.Trigger(eid, "foo")
	wt.Trigger(eid, "bar")

	v := <-ch
	if g, w := fmt.Sprintf("%v (%T)", v, v), "foo (string)"; g != w {
		t.Errorf("<-ch = %v, want %v", g, w)
	}

	if g := <-ch; g != nil {
		t.Errorf("unexpected non-nil value: %v (%T)", g, g)
	}
}

func TestIsRegistered(t *testing.T) {
	wt := newWait()

	wt.Register(0)
	wt.Register(1)
	wt.Register(2)

	for i := uint64(0); i < 3; i++ {
		if !wt.IsRegistered(i) {
			t.Errorf("event ID %d isn't registered", i)
		}
	}

	if wt.IsRegistered(4) {
		t.Errorf("event ID 4 shouldn't be registered")
	}

	wt.Trigger(0, "foo")
	if wt.IsRegistered(0) {
		t.Errorf("event ID 0 is already triggered, shouldn't be registered")
	}
}
