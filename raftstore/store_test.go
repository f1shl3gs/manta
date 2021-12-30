package raftstore

/*
import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/f1shl3gs/manta/raftstore/rawkv"
	"github.com/f1shl3gs/manta/raftstore/transport"
	"github.com/hashicorp/go-sockaddr"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
)

const (
	testAddr = "localhost:8089"
)

func TestStore(t *testing.T) {
	addr, err := sockaddr.GetPrivateIP()
	require.NoError(t, err)
	addr = testAddr

	testDir := "tests"
	cf := NewConfig()
	logger := zaptest.NewLogger(t, zaptest.Level(zapcore.InfoLevel))

	_ = os.RemoveAll(testDir)
	err = os.MkdirAll(testDir, 0777)
	require.NoError(t, err)

	cf.InitialPeers = []string{addr}
	cf.BindAddr = addr
	cf.DataDir = filepath.Join(testDir, cf.DataDir)
	cf.WALDir = filepath.Join(testDir, cf.WALDir)
	cf.SnapDir = filepath.Join(testDir, cf.SnapDir)

	err = os.MkdirAll(cf.DataDir, 0777)
	require.NoError(t, err)

	store, err := New(cf, logger)
	require.NoError(t, err)

	go func() {
		svr := grpc.NewServer()
		transport.RegisterRaftServer(svr, store)
		rawkv.RegisterRawKVServer(svr, store)

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
		err = http.ListenAndServe(":8087", nil)
		if err != nil {
			panic(err)
		}
	}()

	select {}
}

func TestPut(t *testing.T) {
	cc, err := grpc.Dial(testAddr, grpc.WithInsecure())
	require.NoError(t, err)
	defer cc.Close()

	cli := rawkv.NewRawKVClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cli.Put(ctx, &rawkv.PutRequest{
		Key:   []byte("foo"),
		Value: []byte("bar"),
	})
	require.NoError(t, err)
}

func TestGet(t *testing.T) {
	cc, err := grpc.Dial(testAddr, grpc.WithInsecure())
	require.NoError(t, err)
	defer cc.Close()

	cli := rawkv.NewRawKVClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := []byte("foo")

	_, err = cli.Get(ctx, &rawkv.GetRequest{Key: key})
	require.NoError(t, err)
}

func bench(t *testing.T, connN, workerN int, total int64, benchFn func(ctx context.Context, id int64, cli rawkv.RawKVClient)) {
	remain := total
	wg := &sync.WaitGroup{}
	valueSize := 256

	value := make([]byte, 256)
	for i := 0; i < valueSize; i++ {
		value[i] = 'a'
	}

	conns := make([]*grpc.ClientConn, 0, connN)
	// prepare grpc connection
	for i := 0; i < connN; i++ {
		cc, err := grpc.Dial(testAddr, grpc.WithInsecure())
		require.NoError(t, err)

		conns = append(conns, cc)
	}

	workerFn := func(cli rawkv.RawKVClient) {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		for {
			n := atomic.AddInt64(&remain, -1)
			if n <= 0 {
				break
			}

			benchFn(ctx, n, cli)
		}
	}

	defer func() {
		for _, cc := range conns {
			cc.Close()
		}
	}()

	start := time.Now()
	wg.Add(workerN)
	for i := 0; i < workerN; i++ {
		cc := conns[i%connN]

		go workerFn(rawkv.NewRawKVClient(cc))
	}

	wg.Wait()

	elapsed := time.Since(start)
	send := int64(8+len(value)) * total
	rate := float64(total) / elapsed.Seconds()

	fmt.Println("Total", total)
	fmt.Println("Connections", connN)
	fmt.Println("Clients", workerN)
	fmt.Println("Key/Value Size", 8, len(value))
	fmt.Println("Time", elapsed)
	fmt.Println("Throughput", float64(send)/1024.0/1024.0/elapsed.Seconds(), "MB/s")
	fmt.Println("OPS", rate)
}

const (
	valueSize = 256
)

func TestBenchPut(t *testing.T) {
	value := make([]byte, valueSize)
	for i := 0; i < valueSize; i++ {
		value[i] = 'a'
	}

	bench(t, 100, 1000, 100000, func(ctx context.Context, id int64, cli rawkv.RawKVClient) {
		var key = make([]byte, 8)

		binary.BigEndian.PutUint64(key, uint64(id))
		_, err := cli.Put(ctx, &rawkv.PutRequest{
			Key:   key,
			Value: value,
		})
		if err != nil {
			panic(err)
		}
	})
}

func TestBenchGet(t *testing.T) {
	bench(t, 100, 1000, 100000, func(ctx context.Context, id int64, cli rawkv.RawKVClient) {
		var key = make([]byte, 8)

		binary.BigEndian.PutUint64(key, uint64(id))

		resp, err := cli.Get(ctx, &rawkv.GetRequest{
			Key: key,
		})
		if err != nil {
			panic(err)
		}

		if len(resp.Value) != valueSize {
			panic("unexpected value size")
		}
	})
}
*/
