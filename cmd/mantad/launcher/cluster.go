package launcher

import (
	"github.com/f1shl3gs/manta/raftstore"
	"github.com/f1shl3gs/manta/raftstore/transport"

	"go.uber.org/zap"
)

func (l *Launcher) setupCluster(logger *zap.Logger) error {
	cf := raftstore.NewConfig()
	store, err := raftstore.New(cf, logger)
	if err != nil {
		return err
	}

	transport.RegisterRaftServer(l.grpcServer, store)

	return nil
}
