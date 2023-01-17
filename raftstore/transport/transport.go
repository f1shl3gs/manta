package transport

import (
	"errors"
	"strconv"
	"sync"

	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
)

var (
	ErrPeerAlreadyAdded = errors.New("peer already added")
	ErrPeerNotFound     = errors.New("peer not found")
)

type Transporter struct {
	logger *zap.Logger

	mtx   sync.RWMutex
	peers map[uint64]*peer
}

func New(logger *zap.Logger) *Transporter {
	return &Transporter{
		logger: logger,
		peers:  make(map[uint64]*peer),
	}
}

func (t *Transporter) Send(msgs []raftpb.Message) {
	for _, m := range msgs {
		if m.To == 0 {
			return
		}

		t.mtx.RLock()
		peer := t.peers[m.To]
		t.mtx.RUnlock()

		if peer != nil {
			peer.send(m)
			continue
		}

		t.logger.Warn("ignored message send request; unknown remote peer target",
			zap.String("type", m.Type.String()),
			zap.String("peer", strconv.FormatUint(m.To, 16)))
	}
}

func (t *Transporter) Stop() {
	t.mtx.Lock()
	peers := t.peers
	t.peers = nil
	t.mtx.Unlock()

	for _, p := range peers {
		p.stop()
	}
}

// AddPeer adds new peer with id and address to transport.
// If there is already peer with such id in transport, it will
// return error if address is different(in which case UpdatePeer
// should be used) or nil otherwise.
func (t *Transporter) AddPeer(id uint64, addr string) error {
	t.mtx.RLock()
	prev, ok := t.peers[id]
	t.mtx.RUnlock()

	if ok {
		if prev.addr == addr {
			return nil
		}

		return ErrPeerAlreadyAdded
	}

	t.logger.Info("add transport peer",
		zap.String("id", strconv.FormatUint(id, 16)),
		zap.String("addr", addr))

	p, err := newPeer(id, addr)
	if err != nil {
		return err
	}

	t.mtx.Lock()
	t.peers[id] = p
	t.mtx.Unlock()

	return nil
}

// RemovePeer removes peer from transport and wait for it to stop
func (t *Transporter) RemovePeer(id uint64) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	p, ok := t.peers[id]
	if !ok {
		return ErrPeerNotFound
	}

	p.stop()
	delete(t.peers, id)

	return nil
}
