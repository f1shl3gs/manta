package kv_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/kv"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func BenchmarkCreateNode(b *testing.B) {
	store, closer := NewTestBolt(b, false)
	defer closer()

	svc := kv.NewService(zaptest.NewLogger(b), store)
	err := kv.Initial(context.Background(), store)
	require.NoError(b, err)

	org := &manta.Organization{
		Name: "foo",
		Desc: "bar",
	}
	err = svc.CreateOrganization(context.Background(), org)

	orgID := org.ID

	env := "pro"
	status := "on"
	region := "SH01"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		addr := "1.1.1." + strconv.Itoa(i)
		hostname := addr

		node := &manta.Node{
			OrgID:       orgID,
			Address:     addr,
			Hostname:    hostname,
			Env:         env,
			Status:      status,
			Region:      region,
			Annotations: nil,
		}

		err = svc.CreateNode(context.Background(), node)
		if err != nil {
			panic(err)
		}
	}
}

/*
sync with default mmap size
goos: linux
goarch: amd64
pkg: github.com/f1shl3gs/manta/kv
BenchmarkCreateNode_Sync
BenchmarkCreateNode_Sync-32    	     276	   4207457 ns/op	   36858 B/op	     153 allocs/op

sync and 32M mmap size
goos: linux
goarch: amd64
pkg: github.com/f1shl3gs/manta/kv
BenchmarkCreateNode_Sync
BenchmarkCreateNode_Sync-32    	     280	   4112505 ns/op	   36819 B/op	     151 allocs/op

nosync
goos: linux
goarch: amd64
pkg: github.com/f1shl3gs/manta/kv
BenchmarkCreateNode_Sync
BenchmarkCreateNode_Sync-32    	   21588	     53275 ns/op	   56357 B/op	     177 allocs/op
*/
