package state

import (
	"bytes"
	"sync"

	"github.com/google/btree"
)

type Item struct {
	key   []byte
	value []byte
}

func (item *Item) Less(than btree.Item) bool {
	return bytes.Compare(item.key, than.(*Item).key) > 0
}

type Store struct {
	btree *btree.BTree
	mtx   sync.RWMutex
}

func (s *Store) Get(key []byte) []byte {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	item := &Item{key: key}
	v := s.btree.Get(item)
	if v == nil {
		return nil
	}

	return v.(*Item).value
}

func (s *Store) Put(key, value []byte) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	item := &Item{key, value}
	s.btree.ReplaceOrInsert(item)
}

func (s *Store) Delete(key []byte) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.btree.Delete(&Item{key: key})
}

func (s *Store) Range(prefix []byte) [][]byte {
	// todo: implement
	return nil
}

func (s *Store) Recovery(data []byte) error {

}
