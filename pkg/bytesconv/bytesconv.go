package bytesconv

import (
	"reflect"
	"unsafe"
)

func StringToBytes(s string) (b []byte) {
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
