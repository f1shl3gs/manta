package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemove(t *testing.T) {
	slice := []string{"a", "b", "c"}
	require.Equal(t, []string{"b", "c"}, Remove("a", slice))
	require.Equal(t, []string{"a", "c"}, Remove("b", slice))
	require.Equal(t, []string{"a", "b"}, Remove("c", slice))
}
