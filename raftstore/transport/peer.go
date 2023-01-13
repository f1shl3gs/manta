package transport

import (
	"context"
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
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

	cc     *grpc.ClientConn
	client RaftInternalClient

	mtx         sync.RWMutex
	active      bool
	activeSince time.Time
}

func newPeer(id uint64, addr string) (*Peer, error) {
	peer := &Peer{
		id:     id,
		addr:   addr,
		msgCh:  make(chan raftpb.Message, 64),
		stopCh: make(chan struct{}),
	}

	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// send message asynchronously
	go func() {
		defer cc.Close()

		cli := NewRaftInternalClient(cc)

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

func (peer *Peer) send(msg raftpb.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultConnWriteTimeout)
	defer cancel()

	_, err := peer.client.Send(ctx, &msg)
	if err != nil {
		peer.logger.Warn("send message failed",
			zap.Error(err))
		return
	}

	return
}

func (peer *Peer) stop() {
	close(peer.stopCh)
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
