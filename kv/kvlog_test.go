package kv_test

import (
	"crypto/sha1"
	"fmt"
	"testing"
	"time"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/pkg/snowflake"
	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	key := []byte("fooooooo")
	h := sha1.New()
	h.Write(key)
	v := h.Sum(nil)

	fmt.Printf("%x\n", v)
	fmt.Println(len(v))
}

func TestTime(t *testing.T) {
	now := time.Now()
	us := now.UnixNano()

	pn := time.Unix(0, us)

	require.Equal(t, now, pn)
}

func TestMagic(t *testing.T) {
	g := snowflake.New(0)
	id := manta.ID(g.Next())

	fmt.Println(id.String())
}
