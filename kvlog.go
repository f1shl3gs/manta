package manta

import (
	"context"
	"time"
)

type KeyValueLog interface {
	// AddLogEntry adds an entry to the file
	AddLogEntry(ctx context.Context, k, v []byte, t time.Time) error

	// ForEachLogEntry iterator through all the file entries at key and applies the function fn for each record
	ForEachLogEntry(ctx context.Context, k []byte, fn func(v []byte, t time.Time) error) error
}
