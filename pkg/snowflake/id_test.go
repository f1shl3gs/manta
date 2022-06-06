package snowflake

import (
	"fmt"
	"testing"
)

func TestIdGenerator_ID(t *testing.T) {
	g := NewIDGenerator()

	for i := 0; i < 10; i++ {
		id := g.ID()
		fmt.Println(id.String())
	}
}
