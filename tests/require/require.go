package require

import "fmt"

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})

	FailNow()
}

type tHelper interface {
	Helper()
}

func NoError(t TestingT, err error, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if err == nil {
		return
	}

	Fail(t, fmt.Sprintf("Received unexpected error:\n%+v", err), msgAndArgs...)
	t.FailNow()
}
