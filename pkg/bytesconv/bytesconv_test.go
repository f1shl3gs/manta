package bytesconv

import "testing"

func TestModify(t *testing.T) {
	foo := BytesToString([]byte("foo"))

	if foo+"bar" != "foobar" {
		t.Fatalf("add failed")
	}
}
