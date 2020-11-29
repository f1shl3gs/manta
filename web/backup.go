package web

import (
	"net/http"
	"time"

	"github.com/f1shl3gs/manta"
	"go.uber.org/zap"
)

type BackupHandler struct {
	logger *zap.Logger

	backupService manta.BackupService
}

func (b *BackupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	err := b.backupService.Backup(r.Context(), w)
	latency := time.Since(start)

	if err != nil {
		b.logger.Error("backup failed",
			zap.String("latency", latency.String()),
			zap.Error(err))
	} else {
		b.logger.Info("backup success",
			zap.String("latency", latency.String()))
	}
}
