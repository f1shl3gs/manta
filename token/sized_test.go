package token

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGen(t *testing.T) {
	gen := NewGenerator(0)
	token, err := gen.Token()
	require.NoError(t, err)
	require.Equal(t, len(token) == base64.URLEncoding.EncodedLen(defaultTokenSize), true)
	fmt.Println(token)
}

func BenchmarkGen(b *testing.B) {
	gen := NewGenerator(0)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.Token()
		if err != nil {
			panic("generate token failed")
		}
	}
}

/*
goos: linux
goarch: amd64
pkg: github.com/f1shl3gs/manta/token
cpu: AMD Ryzen 9 3950X 16-Core Processor
BenchmarkGen
BenchmarkGen-32    	  514182	      2306 ns/op	     256 B/op	       3 allocs/op
PASS
*/
