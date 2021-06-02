package transport

import (
	"context"
	"sync"
	"time"

	"github.com/f1shl3gs/manta/raftstore/internal"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	// DefaultConnReadTimeout and DefaultConnWriteTimeout are the i/o timeout set on each connection rafthttp pkg creates.
	// A 5 seconds timeout is good enough for recycling bad connections. Or we have to wait for
	// tcp keepalive failing to detect a bad connection, which is at minutes level.
	// For long term streaming connections, rafthttp pkg sends application level linkHeartbeatMessage
	// to keep the connection alive.
	// For short term pipeline connections, the connection MUST be killed to avoid it being
	// put back to http pkg connection pool.
	DefaultConnReadTimeout  = 5 * time.Second
	DefaultConnWriteTimeout = 5 * time.Second

	ConnectionPoolSize = 8
)

type Peer struct {
	id     uint64
	addr   string
	msgCh  chan raftpb.Message
	stopCh chan struct{}
	logger *zap.Logger

	mtx         sync.RWMutex
	active      bool
	activeSince time.Time
}

func (peer *Peer) send(msg raftpb.Message) {

}

func (peer *Peer) sendSnapshot(msg snap.Message) {
	// start a new *grpc
}

func (peer *Peer) stop() {

}

func newPeer(id uint64, addr string) (*Peer, error) {
	peer := &Peer{
		id:     id,
		addr:   addr,
		stopCh: make(chan struct{}),
	}

	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	go func() {
		defer cc.Close()

		cli := NewRaftClient(cc)

		for {
			select {
			case <-peer.stopCh:
				return
			case msg := <-peer.msgCh:
				_, err = cli.Send(context.Background(), &msg)
				if err != nil {
					peer.logger.Warn("send raft message failed",
						zap.String("to", internal.IDToString(id)),
						zap.String("addr", addr))
					peer.setInactive()
				} else {
					peer.setActive()
				}
			}
		}
	}()

	return peer, nil
}

func (peer *Peer) setActive() {
	peer.mtx.Lock()
	defer peer.mtx.Unlock()

	if !peer.active {
		peer.active = true
		peer.activeSince = time.Now()
	}
}

func (peer *Peer) setInactive() {
	peer.mtx.Lock()
	defer peer.mtx.Unlock()

	peer.active = false
	peer.activeSince = time.Time{}
}

func (peer *Peer) activeTime() time.Time {
	peer.mtx.RUnlock()
	defer peer.mtx.RUnlock()

	return peer.activeSince
}
