package mock

import "time"

type Time struct {
	Current time.Time
}

func (t *Time) Now() time.Time {
	return t.Current
}
