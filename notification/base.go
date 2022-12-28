package notification

import (
	"time"

	"github.com/f1shl3gs/manta"
)

type Base struct {
	ID      manta.ID  `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Name string `json:"name"`
	Desc string `json:"desc"`

	OrgID manta.ID `json:"orgID"`
}

func (b *Base) Valid() error {
	if !b.OrgID.Valid() {
		return manta.ErrInvalidOrgID
	}

	return nil
}

func (b *Base) GetID() manta.ID {
	return b.ID
}

func (b *Base) SetID(id manta.ID) {
	b.ID = id
}

func (b *Base) GetOrgID() manta.ID {
	return b.OrgID
}

func (b *Base) SetOrgID(orgID manta.ID) {
	b.OrgID = orgID
}

func (b *Base) GetName() string {
	return b.Name
}

func (b *Base) SetName(name string) {
	b.Name = name
}

func (b *Base) GetDesc() string {
	return b.Desc
}

func (b *Base) SetDesc(desc string) {
	b.Desc = desc
}

// SetCreated implement manta.CRUDSetter
func (b *Base) SetCreated(ts time.Time) {
	b.Created = ts
}

// SetUpdated implement manta.CRUDSetter
func (b *Base) SetUpdated(ts time.Time) {
	b.Updated = ts
}

// GetCreated implement manta.CRUDGetter
func (b *Base) GetCreated() time.Time {
	return b.Created
}

// GetUpdated implement manta.CRUDGetter
func (b *Base) GetUpdated() time.Time {
	return b.Updated
}
