package raftstore

import (
	"context"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"testing"

	"github.com/f1shl3gs/manta/raftstore/transport"
	"github.com/hashicorp/go-sockaddr"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
)

func TestStore(t *testing.T) {
	addr, err := sockaddr.GetPrivateIP()
	require.NoError(t, err)
	addr = "localhost:8087"

	testDir := "tests"
	cf := NewConfig()
	logger := zaptest.NewLogger(t)

	_ = os.RemoveAll(testDir)

	cf.InitialPeers = []string{addr}
	cf.BindAddr = addr
	cf.DataDir = filepath.Join(testDir, cf.DataDir)
	cf.WALDir = filepath.Join(testDir, cf.WALDir)
	cf.SnapDir = filepath.Join(testDir, cf.SnapDir)

	store, err := New(cf, logger)
	require.NoError(t, err)

	go func() {
		svr := grpc.NewServer()
		transport.RegisterRaftServer(svr, store)
		l, err := net.Listen("tcp", addr)
		require.NoError(t, err)

		err = svr.Serve(l)
		require.NoError(t, err)
	}()

	go func() {
		err = store.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err = http.ListenAndServe(":8088", nil)
		if err != nil {
			panic(err)
		}
	}()

	// go func() {
	// 	time.Sleep(time.Second)
	// 	store.raftNode.Campaign(context.Background())
	// }()

	select {}
}
