package transport

import (
	"errors"
	"sync"
	"time"

	"github.com/f1shl3gs/manta/raftstore/internal"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/snap"
	"go.uber.org/zap"
)

var (
	errMemberRemoved    = errors.New("the member has been permanently removed from the cluster")
	errMemberNotFound   = errors.New("member not found")
	ErrPeerAlreadyAdded = errors.New("peer already added")
	ErrPeerNotFound     = errors.New("peer not found")
)

type Resolver interface {
	Resolve(id uint64) (string, error)
}

type Transporter interface {
	// Send sends out the given messages to the remove peers.
	// Each message has a To field, which is an id that maps
	// to an existing peer in the transport. If the id cannot
	// be found in the transport, the message will be ignored.
	Send(msgs []raftpb.Message)

	// SendSnapshot sends out the given snapshot message to
	// a remote peer. The behavior of SendSnapshot is similar to Send
	SendSnapshot(m snap.Message)

	// AddPeer adds a peer with given peer urls into the transport.
	// It is the caller's responsibility to ensure the urls are all valid,
	// or it panics
	AddPeer(id uint64, addr string) error

	// RemovePeer removes the peer with given id
	RemovePeer(id uint64) error

	// RemoveAllPeers removes all existing peers in the transport
	RemoveAllPeers()

	// UpdatePeer updates the peer addr of the the peer with the given id.
	// It is the caller's responsibility to ensure the addr are valid, or
	// it panics
	UpdatePeer(id uint64, addr string)

	ActiveSince(id uint64) time.Time

	// ActivePeers returns the number of active peers
	ActivePeers() int
}

var _ Transporter = &transport{}

// Raft is interface which represents Raft API for transport package
type Raft interface {
	// IsIDRemoved is implement by membership.Cluster
	IsIDRemoved(id uint64) bool

	ReportUnreachable(id uint64)

	ReportSnapshot(id uint64, status raft.SnapshotStatus)
}

type transport struct {
	raft   Raft
	logger *zap.Logger

	mtx   sync.RWMutex
	peers map[uint64]*Peer
}

func (trans *transport) Send(msgs []raftpb.Message) {
	for _, m := range msgs {
		if m.To == 0 {
			// ignore intentionally dropped message
			continue
		}

		trans.logger.Debug("send message",
			zap.String("to", internal.IDToString(m.To)),
			zap.String("from", internal.IDToString(m.From)),
			zap.String("type", m.Type.String()))

		to := m.To
		trans.mtx.RLock()
		peer := trans.peers[to]
		trans.mtx.RUnlock()

		if peer != nil {
			peer.send(m)
			continue
		}

		trans.logger.Warn("ignored message send request; unknown remote peer target",
			zap.String("type", m.Type.String()),
			zap.String("peer", internal.IDToString(to)))
	}
}

func (trans *transport) SendSnapshot(m snap.Message) {
	trans.mtx.RLock()
	defer trans.mtx.RUnlock()

	p := trans.peers[m.To]
	if p == nil {
		m.CloseWithError(errMemberNotFound)
		return
	}

	p.sendSnapshot(m)
}

// AddPeer adds new peer with id and address to transport.
// If there is already peer with such id in transport, it will
// return error if address is different(in which case UpdatePeer
// should be used) or nil otherwise.
func (trans *transport) AddPeer(id uint64, addr string) error {
	trans.mtx.Lock()
	defer trans.mtx.Unlock()

	if prev, ok := trans.peers[id]; ok {
		if prev.addr == addr {
			return nil
		}

		return ErrPeerAlreadyAdded
	}

	trans.logger.Info("add transport peer",
		zap.String("id", internal.IDToString(id)),
		zap.String("addr", addr))

	p, err := newPeer(id, addr)
	if err != nil {
		return err
	}

	trans.peers[id] = p

	return nil
}

// RemovePeer removes peer from transport and wait for it to stop
func (trans *transport) RemovePeer(id uint64) error {
	trans.mtx.Lock()
	defer trans.mtx.Unlock()

	p, ok := trans.peers[id]
	if !ok {
		return ErrPeerNotFound
	}

	p.stop()
	delete(trans.peers, id)

	return nil
}

func (trans *transport) RemoveAllPeers() {
	panic("implement me")
}

func (trans *transport) UpdatePeer(id uint64, addr string) {
	panic("implement me")
}

func (trans *transport) ActiveSince(id uint64) time.Time {
	panic("implement me")
}

func (trans *transport) Stop() {
	panic("implement me")
}

func (trans *transport) ActivePeers() int {
	trans.mtx.RLock()
	defer trans.mtx.RUnlock()

	n := 0
	for _, peer := range trans.peers {
		ts := peer.activeTime()
		if !ts.IsZero() {
			n += 1
		}
	}

	return n
}

func New(raft Raft, logger *zap.Logger) Transporter {
	return &transport{
		raft:   raft,
		logger: logger,
		peers:  make(map[uint64]*Peer),
	}
}
