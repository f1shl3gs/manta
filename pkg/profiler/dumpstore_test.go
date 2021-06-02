package profiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDumpStore_GC(t *testing.T) {
	tests := []struct {
		name       string
		startFiles []string
		sizes      []int64
		maxSize    uint64
		cleaned    []string
	}{
		{
			name:       "no profiles",
			startFiles: []string{},
			sizes:      []int64{},
			maxSize:    10,
			cleaned:    []string{},
		},
		{
			name: "cleanup some files",
			startFiles: []string{
				"memprof.2020-06-11T13_19_19.000.1",
				"memprof.2020-06-12T13_19_19.001.2",
				"memprof.2020-06-13T13_19_19.002.3",
				"memprof.2020-06-14T13_19_19.003.4",
				"memprof.2020-06-15T13_19_19.004.5",
			},
			// The actual sizes for the 5 names above.
			sizes: []int64{
				10, // June 11
				10, // June 12
				10, // June 13
				10, // June 14
				10, // June 15
			},
			// The max size to keep.
			maxSize: 35,
			// Expected files to clean up.
			cleaned: []string{
				"memprof.2020-06-12T13_19_19.001.2",
				"memprof.2020-06-11T13_19_19.000.1",
			},
		},
		{
			name: "below maxSize",
			startFiles: []string{
				"memprof.2020-06-11T13_19_19.000.1",
				"memprof.2020-06-12T13_19_19.001.2",
				"memprof.2020-06-13T13_19_19.002.3",
				"memprof.2020-06-14T13_19_19.003.4",
				"memprof.2020-06-15T13_19_19.004.5",
			},
			// The actual sizes for the 5 names above.
			sizes: []int64{
				10, // June 11
				10, // June 12
				10, // June 13
				10, // June 14
				10, // June 15
			},
			// The max size to keep.
			maxSize: 100,
			// Expected files to clean up.
			cleaned: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path, err := ioutil.TempDir("", "remove")
			require.NoError(t, err)

			defer os.RemoveAll(path)

			createTestdata(t, path, tc.startFiles, tc.sizes)

			cleaned := []string{}
			cleanupFn := func(s string) error {
				cleaned = append(cleaned, s)
				return nil
			}

			ds := dumpStore{
				maxSize:   tc.maxSize,
				dir:       path,
				prefix:    "memprof",
				suffix:    "",
				cleanupFn: cleanupFn,
			}

			err = ds.GC()
			require.NoError(t, err)
			require.EqualValues(t, tc.cleaned, cleaned)
		})
	}
}

func createTestdata(t *testing.T, dir string, files []string, sizes []int64) {
	if len(sizes) > 0 {
		require.Equal(t, len(files), len(sizes))
	}

	for i, name := range files {
		f, err := os.Create(filepath.Join(dir, name))
		require.NoError(t, err)

		_, err = fmt.Fprintf(f, "%*s", sizes[i], " ")
		require.NoError(t, err)

		err = f.Close()
		require.NoError(t, err)
	}
}
