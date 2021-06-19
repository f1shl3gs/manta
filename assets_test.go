package manta

import (
	"compress/gzip"
	"fmt"
	"os"
	"testing"

	"github.com/f1shl3gs/manta/pkg/tarfs"
	"github.com/stretchr/testify/require"
)

func TestWalk(t *testing.T) {
	f, err := Assets.Open("assets.tgz")
	require.NoError(t, err)

	defer f.Close()

	gr, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer gr.Close()

	tfs, err := tarfs.New(gr)
	require.NoError(t, err)

	tfs.Walk(func(name string, fi os.FileInfo) {
		fmt.Println(name, fi.Size())
	})
}
