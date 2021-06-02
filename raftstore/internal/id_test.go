package internal

import (
	"fmt"
	"strconv"
	"testing"
)

func TestIDToString(t *testing.T) {
	fmt.Println(strconv.FormatUint(uint64(1), 16))
}
