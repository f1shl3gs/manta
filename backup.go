package manta

import (
	"context"
	"io"
)

type BackupService interface {
	Backup(ctx context.Context, w io.Writer) error
}
