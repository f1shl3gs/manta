package snowflake

import (
	"math/rand"

	"github.com/f1shl3gs/manta"
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
