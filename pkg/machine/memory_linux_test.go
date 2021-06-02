package machine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTotalMemory(t *testing.T) {
	procDir = "testdata"
	n, err := GetTotalMemory()
	require.NoError(t, err)
	require.Equal(t, n, uint64(32778120*1024))
}
