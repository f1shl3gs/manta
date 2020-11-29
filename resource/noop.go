package resource

import (
	"github.com/f1shl3gs/manta"
)

type noop struct{}

func NewNoopResourceLogger() Logger {
	return &noop{}
}

func (n *noop) Log(c manta.Change) error {
	return nil
}
