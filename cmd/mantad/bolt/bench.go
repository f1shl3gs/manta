package bolt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	// ErrInvalidValue is returned when a benchmark reads an unexpected value.
	ErrInvalidValue = errors.New("invalid value")
)

func bench() *cobra.Command {
	var (
		path      string
		noSync    bool
		batchSize int
		count     int
	)

	cmd := &cobra.Command{
		Use:   "bench",
		Short: "bench boltdb",
		RunE: func(cmd *cobra.Command, args []string) error {
			if path == "" {
				f, err := ioutil.TempFile("", "bolt-bench-")
				if err != nil {
					return err
				}

				f.Close()
				os.Remove(f.Name())
				path = f.Name()
			}

			db, err := bolt.Open(path, 0666, nil)
			if err != nil {
				return err
			}
			db.NoSync = noSync
			defer db.Close()

			start := time.Now()
			if err = runWrites(db, count, batchSize); err != nil {
				return err
			}
			duration := time.Since(start)

			fmt.Println("Write:")
			fmt.Printf("    Keys:        %d\n", count)
			fmt.Printf("    Batch-Size:  %d\n", batchSize)
			fmt.Printf("    Duration:    %fs\n", duration.Seconds())
			fmt.Printf("    Throughput:  %f w/s\n", float64(count)/duration.Seconds())

			start = time.Now()
			if reads, err := runReads(db); err != nil {
				return err
			} else {
				duration := time.Since(start)

				fmt.Println("Read:")
				fmt.Printf("    Keys:        %d\n", reads)
				fmt.Printf("    Duration:    %fs\n", duration.Seconds())
				fmt.Printf("    Throughput:  %f r/s\n", float64(count)/duration.Seconds())
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&path, "path", "", "bolt file")
	cmd.Flags().BoolVar(&noSync, "no-sync", false, "sync on every write")
	cmd.Flags().IntVar(&batchSize, "batch-size", 0, "")
	cmd.Flags().IntVar(&count, "count", 1000, "")

	return cmd
}

var benchBucketName = []byte("bench")

func runWrites(db *bolt.DB, count, batchSize int) error {
	keySize := 8
	valueSize := 256
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	keySource := func() uint32 {
		return r.Uint32()
	}

	if batchSize == 0 {
		batchSize = count
	} else if count%batchSize != 0 {
		panic("")
	}

	for i := 0; i < count; i += batchSize {
		if err := db.Update(func(tx *bolt.Tx) error {
			b, _ := tx.CreateBucketIfNotExists(benchBucketName)
			b.FillPercent = 0.5

			for j := 0; j < batchSize; j++ {
				key := make([]byte, keySize)
				value := make([]byte, valueSize)

				// Write key as uint32
				binary.BigEndian.PutUint32(key, keySource())

				// Insert key/value
				if err := b.Put(key, value); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// Read from the database
func runReads(db *bolt.DB) (int, error) {
	var count int

	err := db.View(func(tx *bolt.Tx) error {
		t := time.Now()

		for {
			var top = tx.Bucket(benchBucketName)
			if err := top.ForEach(func(name, _ []byte) error {
				if b := top.Bucket(name); b != nil {
					c := b.Cursor()
					for k, v := c.First(); k != nil; k, v = c.Next() {
						if v == nil {
							return ErrInvalidValue
						}

						count += 1
					}
				}

				return nil
			}); err != nil {
				return err
			}

			// Make sure we do this for at least a second
			if time.Since(t) >= time.Second {
				break
			}
		}

		return nil
	})

	return count, err
}
