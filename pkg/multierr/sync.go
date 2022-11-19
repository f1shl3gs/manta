package multierr

import (
	"sync"
)

type SyncedErrors struct {
	mtx    sync.Mutex
	errors []error
}

func (s *SyncedErrors) Add(err error) {
	if err == nil {
		return
	}

	s.mtx.Lock()
	s.errors = append(s.errors, err)
	s.mtx.Unlock()
}

func (s *SyncedErrors) Unwrap() error {
	if s == nil {
		return nil
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if len(s.errors) == 0 {
		return nil
	}

	var c chain
	for _, e := range s.errors {
		c = append(c, e)
	}

	return c
}
