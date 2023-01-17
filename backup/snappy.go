package backup

import (
	"context"
	"io"

	"github.com/f1shl3gs/manta"

	"github.com/golang/snappy"
)

type Snappy struct {
	service manta.BackupService
}

// Backup implmemt BackupService
func (s *Snappy) Backup(ctx context.Context, w io.Writer) error {
	return s.service.Backup(ctx, snappy.NewBufferedWriter(w))
}
