package store

import (
	"bytes"

	"github.com/google/btree"
)

type item struct {
	key         []byte
	version     uint64
	generations [][]byte
}

func (i *item) Less(than btree.Item) bool {
	t := than.(*item)
	return bytes.Compare(i.key, t.key) > 0
}

type store struct {
	btree btree.BTree
}

func (s *store) scan(key []byte) {
	s.btree.Get(&item{key: key})
}
