package mock

import (
	"github.com/f1shl3gs/manta"
)

type StaticIDGenerator struct {
	id manta.ID
}

func NewStaticIDGenerator(id manta.ID) *StaticIDGenerator {
	return &StaticIDGenerator{id: id}
}

func (s *StaticIDGenerator) ID() manta.ID {
	return s.id
}

func (s *StaticIDGenerator) Set(id manta.ID) {
	s.id = id
}

type IncrementalIDGenerator struct {
	n uint64
}

func (i *IncrementalIDGenerator) ID() manta.ID {
	n := i.n
	i.n++

	return manta.ID(n)
}

func NewIncrementalIDGenerator(start uint64) *IncrementalIDGenerator {
	if start == 0 {
		panic("start id cannot be zero")
	}

	return &IncrementalIDGenerator{
		n: start,
	}
}
