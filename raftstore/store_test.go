package raftstore

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/f1shl3gs/manta/kv"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func setupStore(t testing.TB) *Store {
	l, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	assert.NoError(t, l.Close())
	addr := l.Addr().String()

	cf := &Config{
		Listen:  addr,
		DataDir: t.TempDir(),
	}
	logger := zaptest.NewLogger(t, zaptest.Level(zapcore.InfoLevel))

	store, err := New(cf, logger)
	assert.NoError(t, err)

	return store
}

func TestNew(t *testing.T) {
	go func() {
		_ = http.ListenAndServe(":7000", nil)
	}()

	l, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	assert.NoError(t, l.Close())
	addr := l.Addr().String()

	cf := &Config{
		Listen:  addr,
		DataDir: t.TempDir(),
	}
	logger := zaptest.NewLogger(t)

	store, err := New(cf, logger)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go store.Run(ctx)

	bucketName := []byte("foo")
	key, value := []byte("key"), []byte("value")

	err = store.CreateBucket(ctx, bucketName)
	assert.NoError(t, err)

	err = store.Update(ctx, func(tx kv.Tx) error {
		b, err := tx.Bucket(bucketName)
		if err != nil {
			return err
		}

		return b.Put(key, value)
	})
	assert.NoError(t, err)
	dump(t, store.db.Load())

	// delete key and write another
	err = store.Update(ctx, func(tx kv.Tx) error {
		b, err := tx.Bucket(bucketName)
		if err != nil {
			return err
		}

		if err = b.Delete(key); err != nil {
			return err
		}

		return b.Put([]byte("foooo"), []byte("baaar"))
	})
	dump(t, store.db.Load())
}

func TestRestart(t *testing.T) {
	go func() {
		_ = http.ListenAndServe(":7000", nil)
	}()

	l, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	assert.NoError(t, l.Close())
	addr := l.Addr().String()

	cf := &Config{
		Listen:  addr,
		DataDir: t.TempDir(),
	}
	logger := zaptest.NewLogger(t)

	store, err := New(cf, logger)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doneCh := make(chan struct{})
	go func() {
		store.Run(ctx)
		close(doneCh)
	}()

	bucketName := []byte("foo")
	key, value := []byte("key"), []byte("value")

	err = store.CreateBucket(ctx, bucketName)
	assert.NoError(t, err)

	err = store.Update(ctx, func(tx kv.Tx) error {
		b, err := tx.Bucket(bucketName)
		if err != nil {
			return err
		}

		return b.Put(key, value)
	})
	assert.NoError(t, err)
	dump(t, store.db.Load())

	store.stop()

	fmt.Println("store stopped")

	<-doneCh

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	fmt.Println("restart node")
	store, err = New(cf, logger)
	assert.NoError(t, err)

	go store.Run(ctx)

	select {
	case <-store.ReadyNotify():
	case <-ctx.Done():
		panic("timeout")
	}

	err = store.Update(ctx, func(tx kv.Tx) error {
		b, err := tx.Bucket(bucketName)
		if err != nil {
			return err
		}

		return b.Put([]byte("key1"), []byte("value1"))
	})
	assert.NoError(t, err)
	dump(t, store.db.Load())
}

func dump(t *testing.T, db *bolt.DB) {
	err := db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Println("--------------------------------------------------------------------------------------")
			fmt.Println(string(name))
			return b.ForEach(func(k, v []byte) error {
				fmt.Printf("  %s %s\n", k, v)
				return nil
			})
		})
	})
	fmt.Println()
	fmt.Println()
	assert.NoError(t, err)
}

func TestBench(t *testing.T) {
	t.SkipNow()

	go func() {
		_ = http.ListenAndServe(":7000", nil)
	}()

	ctx := context.Background()
	bucketName := []byte("bn")
	value, err := os.ReadFile("../example/dashboards/manta.json")
	assert.NoError(t, err)

	store := setupStore(t)
	go store.Run(context.Background())

	err = store.CreateBucket(ctx, bucketName)
	assert.NoError(t, err)

	var (
		total   = uint64(50000)
		threads = 32
		counter atomic.Uint64
		wg      sync.WaitGroup
	)

	writer := func() {
		for {
			n := counter.Add(1)
			if n >= total {
				return
			}

			key := unsafeStringToBytes(strconv.FormatUint(n, 16))
			err = store.Update(ctx, func(tx kv.Tx) error {
				b, err := tx.Bucket(bucketName)
				if err != nil {
					return err
				}

				return b.Put(key, value)
			})
			if err != nil {
				panic(err)
			}
		}
	}

	start := time.Now()
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()
			writer()
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	store.sync()
	file := store.db.Load().Path()
	stat, err := os.Stat(file)
	assert.NoError(t, err)

	fmt.Println(file)

	fmt.Printf("Total writes    %d\n", total)
	fmt.Printf("Time:           %s\n", elapsed.String())
	fmt.Printf("Writers:        %d\n", threads)
	fmt.Printf("Throughtput:    %f MB/s\n", float64(total)*float64(len(value))/(elapsed.Seconds()*1024.0*1024.0))
	fmt.Printf("QPS:            %f\n", float64(total)/elapsed.Seconds())
	fmt.Printf("DB Size:        %f MB\n", float64(stat.Size())/1024.0/1024.0)
}
