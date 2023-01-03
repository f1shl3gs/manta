package mvcc

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrRevisionNotFound = errors.New("mvcc revision not found")

	ErrPutWithSmallerRevision = errors.New("put with an unexpected smaller revision")
)

// generation contains multiple revisions of a key
type generation struct {
	ver int64
	// when the generation is created (put in first revision).
	created revision
	revs    []revision
}

func (g *generation) isEmpty() bool {
	return g == nil || len(g.revs) == 0
}

// walk walks through the revisions in the generation in descending order.
// It passes the revision to the given function.
// walk returns until:
//  1. it finishes walking all pairs
//  2. the function returns false.
//
// walk returns the position at where it stopped. If it stopped after
// finishing walking, -1 will be returned.
func (g *generation) walk(fn func(rev revision) bool) int {
	length := len(g.revs)
	for i := range g.revs {
		ok := fn(g.revs[length-i-1])
		if !ok {
			return length - i - 1
		}
	}

	return -1
}

func (g *generation) String() string {
	return fmt.Sprintf("g: created[%d] ver[%d], revs %#v\n", g.created, g.ver, g.revs)
}

func (g generation) equal(b generation) bool {
	if g.ver != b.ver {
		return false
	}
	if len(g.revs) != len(b.revs) {
		return false
	}

	for i := range g.revs {
		ar, br := g.revs[i], b.revs[i]
		if ar != br {
			return false
		}
	}
	return true
}

// keyIndex stores the revisions of a key in the backend.
// Each keyIndex has at least one key generation.
// Each generation might have several key versions.
// Tombstone on a key appends an tombstone version at the end
// of the current generation and creates a new empty generation.
// Each version of a key has an index pointing to the backend.
//
// For example: put(1.0);put(2.0);tombstone(3.0);put(4.0);tombstone(5.0) on key "foo"
// generate a keyIndex:
// key:     "foo"
// modified: 5
// generations:
//
//	{empty}
//	{4.0, 5.0(t)}
//	{1.0, 2.0, 3.0(t)}
//
// Compact a keyIndex removes the versions with smaller or equal to
// rev except the largest one. If the generation becomes empty
// during compaction, it will be removed. if all the generations get
// removed, the keyIndex should be removed.
//
// For example:
// compact(2) on the previous example
// generations:
//
//	{empty}
//	{4.0, 5.0(t)}
//	{2.0, 3.0(t)}
//
// compact(4)
// generations:
//
//	{empty}
//	{4.0, 5.0(t)}
//
// compact(5):
// generations:
//
//	{empty} -> key SHOULD be removed.
//
// compact(6):
// generations:
//
//	{empty} -> key SHOULD be removed.

type keyIndex struct {
	key         []byte
	modified    revision
	generations []generation
}

func (ki *keyIndex) Less(other *keyIndex) bool {
	return bytes.Compare(ki.key, other.key) == -1
}

// put puts a revision to the keyIndex
func (ki *keyIndex) put(main, sub int64) {
	rev := revision{main, sub}

	if !rev.GreaterThan(ki.modified) {
		panic("put with an unexpected smaller revision")
	}

	if len(ki.generations) == 0 {
		ki.generations = append(ki.generations, generation{})
	}

	g := &ki.generations[len(ki.generations)-1]
	if len(g.revs) == 0 {
		// create a new key
		g.created = rev
	}
	g.revs = append(g.revs, rev)
	g.ver++
	ki.modified = rev
}

func (ki *keyIndex) restore(created, modified revision, ver int64) {
	if len(ki.generations) != 0 {
		panic(fmt.Sprintf("'restore' got an unexpected non-empty generations, generations-size: %d", len(ki.generations)))
	}

	ki.modified = modified
	g := generation{
		created: created,
		ver:     ver,
		revs:    []revision{modified},
	}
	ki.generations = append(ki.generations, g)
}

// get gets the modified, created revision and version of the key that satisfies the given atRev.
// Rev must be smaller than or equal to the given atRev
func (ki *keyIndex) get(atRev int64) (modified, created revision, ver int64, err error) {
	if ki.isEmpty() {
		panic("'get' got an unexpected empty keyIndex")
	}

	g := ki.findGeneration(atRev)
	if g.isEmpty() {
		return revision{}, revision{}, 0, ErrRevisionNotFound
	}

	n := g.walk(func(rev revision) bool { return rev.main > atRev })
	if n != -1 {
		return g.revs[n], g.created, g.ver - int64(len(g.revs)-n-1), nil
	}

	return revision{}, revision{}, 0, ErrRevisionNotFound
}

func (ki *keyIndex) isEmpty() bool {
	return len(ki.generations) == 1 && ki.generations[0].isEmpty()
}

// findGeneration finds out the generation of the keyIndex that the
// given rev belongs to. If the given rev is at the gap of two generations,
// which means that the key does not exist at the given rev, it returns nil
func (ki *keyIndex) findGeneration(rev int64) *generation {
	lastg := len(ki.generations) - 1
	cg := lastg

	for cg >= 0 {
		if len(ki.generations[cg].revs) == 0 {
			cg--
			continue
		}

		g := ki.generations[cg]
		if cg != lastg {
			if tomb := g.revs[len(g.revs)-1].main; tomb <= rev {
				return nil
			}
		}

		if g.revs[0].main <= rev {
			return &ki.generations[cg]
		}
		cg--
	}

	return nil
}

// tombstone puts a revision, pointing to a tombstone, to the keyIndex.
// It also creates a new empty generation in the keyIndex.
// It returns ErrRevisionNotFound when tombstone on an empty generation.
func (ki *keyIndex) tombstone(main, sub int64) error {
	if ki.isEmpty() {
		panic(fmt.Sprintf("'tombstone' got an unexpected empty keyIndex, key: %s", string(ki.key)))
	}

	if ki.generations[len(ki.generations)-1].isEmpty() {
		return ErrRevisionNotFound
	}

	ki.put(main, sub)
	ki.generations = append(ki.generations, generation{})
	return nil
}
