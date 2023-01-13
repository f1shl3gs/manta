package mmap

// ZeroOut zeroes out all the bytes in the range [start, end).
func ZeroOut(dst []byte, start, end int) {
	if start < 0 || start >= len(dst) {
		return // BAD
	}
	if end >= len(dst) {
		end = len(dst)
	}
	if end-start <= 0 {
		return
	}
	Memclr(dst[start:end])
	// b := dst[start:end]
	// for i := range b {
	// 	b[i] = 0x0
	// }
}
