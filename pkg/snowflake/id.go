package snowflake

import (
	"github.com/f1shl3gs/manta"
	"math/rand"
)

type IDGenerator struct {
	gen *Generator
}

func (i *IDGenerator) ID() manta.ID {
	for {
		id := manta.ID(i.gen.Next())

		if id.Valid() {
			return id
		}
	}
}

func NewIDGenerator() *IDGenerator {
	idGen := New(rand.Intn(1023))

	return &IDGenerator{
		idGen,
	}
}
