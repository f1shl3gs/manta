package transport

import (
	"container/list"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	IDLETimeout = 5 * time.Minute
	MaxPoolSize = 16
)

type pool struct {
	addr string

	mtx     sync.Mutex
	clients list.List
}

func (pool *pool) Get() (RaftClient, error) {
	var (
		cli RaftClient
	)

	pool.mtx.Lock()
	elmt := pool.clients.Front()
	size := pool.clients.Len()
	pool.mtx.Unlock()

	if elmt == nil {
		// create new connection
		cc, err := grpc.Dial(pool.addr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		cli = NewRaftClient(cc)
	} else {
		cli = elmt.Value.(RaftClient)
	}

	// try to gc
	if size > MaxPoolSize {
		pool.mtx.Lock()
		elmt = pool.clients.Back()
		pool.mtx.Unlock()

		idle := elmt.Value.(*raftClient)
		_ = idle.cc.Close()
	}

	return cli, nil
}

func (pool *pool) Put(cli RaftClient) {
	pool.mtx.Lock()
	pool.clients.PushFront(cli)
	pool.mtx.Unlock()
}
