package broadcast

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	t.Run("no sub", func(t *testing.T) {
		b := New[int]()

		b.Pub(2)
	})

	t.Run("Sub Pub", func(t *testing.T) {
		var (
			b     = New[int]()
			wg    sync.WaitGroup
			got   int
			input = 124
		)

		wg.Add(1)
		go func() {
			defer wg.Done()

			q := b.Sub()
			got = <-q.Ch()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(time.Second)
			b.Pub(input)
		}()

		wg.Wait()

		assert.Equal(t, got, input)
	})

	t.Run("close", func(t *testing.T) {
		var (
			b      = New[int]()
			wg     sync.WaitGroup
			closed bool
		)

		wg.Add(1)
		go func() {
			defer wg.Done()

			q := b.Sub()

			<-q.Ch()
			closed = true
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			time.Sleep(time.Second)

			b.Close()
		}()

		wg.Wait()

		assert.Equal(t, closed, true)
	})
}
