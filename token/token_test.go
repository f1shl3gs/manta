package token

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGen(t *testing.T) {
	gen := NewTokenGenerator(0)
	token, err := gen.Token()
	require.NoError(t, err)
	require.Equal(t, len(token) == base64.URLEncoding.EncodedLen(defaultTokenSize), true)
}

func BenchmarkGen(b *testing.B) {
	gen := NewTokenGenerator(0)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.Token()
		if err != nil {
			panic(err)
		}
	}
}
