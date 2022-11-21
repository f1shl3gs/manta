package duration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	d, err := Parse("2d")
	require.NoError(t, err)
	require.Equal(t, 2*24*time.Hour, d)

	w, err := Parse("3w")
	require.NoError(t, err)
	require.Equal(t, 3*7*24*time.Hour, w)

	s, err := Parse("300s")
	require.NoError(t, err)
	require.Equal(t, 300*time.Second, s)
}
