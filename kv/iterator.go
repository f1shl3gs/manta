package kv

import (
	"context"

	"github.com/f1shl3gs/manta/pkg/tracing"
)

type FilterFn func(key []byte, decodedVal interface{}) bool

type iterator struct {
	cursor Cursor

	counter    int
	descending bool
	limit      int
	offset     int
	prefix     []byte

	nextFn func() (key, val []byte)

	decodeFn func(key, val []byte) (k []byte, decodedVal interface{}, err error)
	filterFn FilterFn
}

func (i *iterator) Next(ctx context.Context) (key []byte, val interface{}, err error) {
	span, _ := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	if i.limit > 0 && i.counter >= i.limit+i.offset {
		return nil, nil, nil
	}

	var k, vRaw []byte
	switch {
	case i.nextFn != nil:
		k, vRaw = i.nextFn()
	case len(i.prefix) > 0:
		k, vRaw = i.cursor.Seek(i.prefix)
		i.nextFn = i.cursor.Next
	case i.descending:
		k, vRaw = i.cursor.Last()
		i.nextFn = i.cursor.Prev
	default:
		k, vRaw = i.cursor.First()
		i.nextFn = i.cursor.Next
	}

	k, decodedVal, err := i.decodeFn(k, vRaw)
	for ; ; k, decodedVal, err = i.decodeFn(i.nextFn()) {
		if err != nil {
			return nil, nil, err
		}
		if i.isNext(k, decodedVal) {
			break
		}
	}
	return k, decodedVal, nil
}

func (i *iterator) isNext(k []byte, v interface{}) bool {
	if len(k) == 0 {
		return true
	}

	if i.filterFn != nil && !i.filterFn(k, v) {
		return false
	}

	// increase counter here since the entity is a valid ent
	// and counts towards the total the user is looking for
	// 	i.e. limit = 5 => 5 valid ents
	//	i.e. offset = 5 => return valid ents after seeing 5 valid ents
	i.counter++

	if i.limit > 0 && i.counter >= i.limit+i.offset {
		return true
	}
	if i.offset > 0 && i.counter <= i.offset {
		return false
	}
	return true
}
