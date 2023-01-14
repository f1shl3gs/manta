package bytesconv

import (
	"unsafe"
)

// StringToBytes performs unholy acts to avoid allocations
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// BytesToString performs unholy acts to avoid allocations
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
