package transport

import (
	"testing"
	"time"
)

func BenchmarkTimeNow(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/f1shl3gs/manta/raftstore/transport
cpu: AMD Ryzen 9 3950X 16-Core Processor
BenchmarkTimeNow
BenchmarkTimeNow-32    	23970580	        48.86 ns/op	       0 B/op	       0 allocs/op
*/

func BenchmarkTimestamp(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Timestamp()
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/f1shl3gs/manta/raftstore/transport
cpu: AMD Ryzen 9 3950X 16-Core Processor
BenchmarkTimestamp
BenchmarkTimestamp-32    	56471276	        20.94 ns/op	       0 B/op	       0 allocs/op
*/
