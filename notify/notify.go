package notify

import (
	"context"

	"github.com/f1shl3gs/manta"
)

type Notifier interface {
	Notify(ctx context.Context, a *manta.Alert)
}
