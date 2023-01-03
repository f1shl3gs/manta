package mvcc

import (
	"sync"

	"github.com/google/btree"
)

type treeIndex struct {
	sync.RWMutex
	tree *btree.BTreeG[*keyIndex]
}

func newTreeIndex() *treeIndex {
	return &treeIndex{
		tree: btree.NewG(32, func(aki *keyIndex, bki *keyIndex) bool {
			return aki.Less(bki)
		}),
	}
}

func (index *treeIndex) Put(key []byte, rev revision) {
	ki := &keyIndex{key: key}

	index.Lock()
	defer index.Unlock()

	old, ok := index.tree.Get(ki)
	if !ok {
		ki.put(rev.main, rev.sub)
		index.tree.ReplaceOrInsert(ki)
	} else {
		old.put(rev.main, rev.sub)
	}
}

func (index *treeIndex) Get(key []byte, atRev int64) (modified, created revision, ver int64, err error) {
	index.RLock()
	defer index.RUnlock()

	return index.unsafeGet(key, atRev)
}

func (index *treeIndex) unsafeGet(key []byte, atRev int64) (modified, created revision, ver int64, err error) {
	ki := &keyIndex{key: key}
	if ki = index.keyIndex(ki); ki == nil {
		return revision{}, revision{}, 0, ErrRevisionNotFound
	}

	return ki.get(atRev)
}

func (index *treeIndex) keyIndex(ki *keyIndex) *keyIndex {
	if ki, ok := index.tree.Get(ki); ok {
		return ki
	}

	return nil
}

func (index *treeIndex) insert(ki *keyIndex) {
	index.Lock()
	defer index.Unlock()

	index.tree.ReplaceOrInsert(ki)
}
