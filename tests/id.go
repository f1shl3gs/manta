package tests

import "github.com/f1shl3gs/manta"

// IDPtr returns a pointer to an manta ID.
func IDPtr(id manta.ID) *manta.ID {
	return &id
}
