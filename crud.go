package manta

import "time"

type CRUDSetter interface {
	SetCreated(ts time.Time)
	SetUpdated(ts time.Time)
}

type CRUDGetter interface {
	GetCreated() time.Time
	GetUpdated() time.Time
}
