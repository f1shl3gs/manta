package transport

import (
	"go.etcd.io/raft/v3/raftpb"
	"go.uber.org/zap"
	"sync"
)

type Transporter struct {
	logger *zap.Logger

	mtx   sync.RWMutex
	peers map[uint64]*Peer
}

func (t *Transporter) Send(msgs []raftpb.Message) {
	/*    batches := make(map[uint64][]raftpb.Message)
	    for _, m := range msgs {
	        if m.To == 0 {
	            return
	        }

	        batches[m.To] = append(batches[m.To], m)
	    }

		for to, batch := range batches {
	        t.mtx.RLock()
	        peer := t.peers[to]
	        t.mtx.RUnlock()

	        if peer != nil {

	        }
	    }*/

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

		t.logger.Debug("ignored message send request; unknown remote peer target",
			zap.String("type", m.Type.String()),
			zap.String("peer", m.To))
	}
}
