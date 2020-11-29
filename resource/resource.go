package resource

import "github.com/f1shl3gs/manta"

type Logger interface {
	Log(c manta.Change) error
}
