package mock

import (
	"github.com/f1shl3gs/manta"
)

type IDGenerator struct {
	Next manta.ID
}

func (m *IDGenerator) ID() manta.ID {
	return m.Next
}
