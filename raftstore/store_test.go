package raftstore

import (
	"context"
	"fmt"
    "github.com/f1shl3gs/manta"
    "net"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"

	"github.com/f1shl3gs/manta/kv"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
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
	logger := zaptest.NewLogger(t)

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

func BenchmarkOneNodeWrite(b *testing.B) {
    bucketName := []byte("bn")
    value := []byte(`{
"id": "0a66228cdb616000",
"created": "2022-12-07T01:47:33.101807828+08:00",
"updated": "2022-12-07T01:57:01.526551906+08:00",
"name": "Manta",
"desc": "Metrics of Manta",
"orgID": "0a659bccc2aba000",
"cells": [
{
"id": "0a659dc98fe16000",
"name": "CPU",
"x": 0,
"y": 4,
"w": 4,
"h": 4,
"viewProperties": {
"type": "xy",
"axes": {
"x": {
"base": ""
},
"y": {
"suffix": "%",
"base": ""
}
},
"queries": [
{
"name": "query 1",
"text": "rate(process_cpu_seconds_total[1m]) * 100"
}
],
"timeFormat": "HH:mm:ss",
"xColumn": "_time",
"yColumn": "_value",
"hoverDimension": "x",
"position": "overlaid",
"geom": "line",
"colors": []
}
},
{
"id": "0a65b658d8216000",
"name": "Threads",
"x": 8,
"y": 4,
"w": 4,
"h": 4,
"viewProperties": {
"type": "xy",
"axes": {
"x": {
"base": ""
},
"y": {
"base": ""
}
},
"queries": [
{
"name": "query 1",
"text": "go_threads"
}
],
"xColumn": "_time",
"yColumn": "_value",
"hoverDimension": "auto",
"position": "overlaid",
"geom": "line",
"colors": []
}
},
{
"id": "0a65effc9de16000",
"name": "Local Object Store Reads",
"x": 0,
"y": 8,
"w": 6,
"h": 4,
"viewProperties": {
"type": "xy",
"axes": {
"x": {
"base": ""
},
"y": {
"base": ""
}
},
"queries": [
{
"name": "query 1",
"text": "increase(boltdb_reads_total[1m])"
}
],
"xColumn": "_time",
"yColumn": "_value",
"hoverDimension": "auto",
"position": "overlaid",
"geom": "line",
"colors": []
}
},
{
"id": "0a65f0208f216000",
"name": "Local Object Store Writes",
"x": 6,
"y": 8,
"w": 6,
"h": 4,
"viewProperties": {
"type": "xy",
"axes": {
"x": {
"base": ""
},
"y": {
"base": ""
}
},
"queries": [
{
"name": "query 1",
"text": "increase(boltdb_writes_total[1m])"
}
],
"xColumn": "_time",
"yColumn": "_value",
"hoverDimension": "auto",
"position": "overlaid",
"geom": "line",
"colors": []
}
},
{
"id": "0a65f0585a216000",
"name": "Memory",
"x": 4,
"y": 4,
"w": 4,
"h": 4,
"viewProperties": {
"type": "xy",
"axes": {
"x": {
"base": ""
},
"y": {
"base": "2"
}
},
"queries": [
{
"name": "query 1",
"text": "process_resident_memory_bytes"
}
],
"xColumn": "_time",
"yColumn": "_value",
"hoverDimension": "auto",
"position": "overlaid",
"geom": "line"
}
},
{
"id": "0a66206d84616000",
"name": "Orgs",
"x": 0,
"y": 0,
"w": 3,
"h": 2,
"viewProperties": {
"type": "single-stat",
"queries": [
{
"name": "query 1",
"text": "boltdb_keys_total{bucket=\"organizations\"}"
}
],
"colors": [
{
"id": "base",
"type": "text",
"hex": "#00C9FF",
"name": "laser"
}
]
}
},
{
"id": "0a6620f64a216000",
"name": "Scrapes",
"x": 6,
"y": 2,
"w": 3,
"h": 2,
"viewProperties": {
"type": "single-stat",
"queries": [
{
"name": "query 1",
"text": "boltdb_keys_total{bucket=\"scrapes\"}"
}
],
"colors": [
{
"id": "base",
"type": "text",
"hex": "#00C9FF",
"name": "laser"
}
]
}
},
{
"id": "0a6622a280e16000",
"name": "Checks",
"x": 9,
"y": 2,
"w": 3,
"h": 2,
"viewProperties": {
"type": "single-stat",
"queries": [
{
"name": "query 1",
"text": "boltdb_keys_total{bucket=\"checks\"}"
}
],
"colors": [
{
"id": "base",
"type": "text",
"hex": "#00C9FF",
"name": "laser"
}
]
}
},
{
"id": "0a6622d8dd216000",
"name": "Users",
"x": 3,
"y": 0,
"w": 3,
"h": 2,
"viewProperties": {
"type": "single-stat",
"queries": [
{
"name": "query 1",
"text": "boltdb_keys_total{bucket=\"users\"}"
}
],
"colors": [
{
"id": "base",
"type": "text",
"hex": "#00C9FF",
"name": "laser"
}
]
}
},
{
"id": "0a6623e12f216000",
"name": "Dashboards",
"x": 6,
"y": 0,
"w": 3,
"h": 2,
"viewProperties": {
"type": "single-stat",
"queries": [
{
"name": "query 1",
"text": "boltdb_keys_total{bucket=\"dashboards\"}"
}
],
"colors": [
{
"id": "base",
"type": "text",
"hex": "#00C9FF",
"name": "laser"
}
]
}
},
{
"id": "0a6623f78aa16000",
"name": "Configs",
"x": 9,
"y": 0,
"w": 3,
"h": 2,
"viewProperties": {
"type": "single-stat",
"queries": [
{
"name": "query 1",
"text": "boltdb_keys_total{bucket=\"configurations\"}"
}
],
"colors": [
{
"id": "base",
"type": "text",
"hex": "#00C9FF",
"name": "laser"
}
]
}
},
{
"id": "0a66245b95216000",
"name": "Task runs",
"x": 3,
"y": 2,
"w": 3,
"h": 2,
"viewProperties": {
"type": "single-stat",
"queries": [
{
"name": "query 1",
"text": "task_scheduler_total_schedule_calls"
}
],
"colors": [
{
"id": "base",
"type": "text",
"hex": "#00C9FF",
"name": "laser"
}
]
}
},
{
"id": "0a66249aeda16000",
"name": "Sessions",
"x": 0,
"y": 2,
"w": 3,
"h": 2,
"viewProperties": {
"type": "single-stat",
"queries": [
{
"name": "query 1",
"text": "boltdb_keys_total{bucket=\"sessions\"}"
}
],
"colors": [
{
"id": "base",
"type": "text",
"hex": "#00C9FF",
"name": "laser"
}
]
}
}
]
}
`)
    ctx := context.Background()

    store := setupStore(b)
    go store.Run(context.Background())

    err := store.CreateBucket(ctx, bucketName)
    assert.NoError(b, err)

    idgen := newGenerator(1, time.Now())
    b.ResetTimer()
	for i := 0; i < b.N; i++ {
        err = store.Update(ctx, func(tx kv.Tx) error {
            b, err := tx.Bucket(bucketName)
            if err != nil {
                return err
            }

            key, err := manta.ID(idgen.Next()).Encode()
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
