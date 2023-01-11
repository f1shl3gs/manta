package transport

import (
	"syscall"
	"time"
)

// Timestamp returns the unix second, the result is effected by
// NTP
func Timestamp() int64 {
	var tv syscall.Timeval

	if err := syscall.Gettimeofday(&tv); err != nil {
		return time.Now().Unix()
	}

	return tv.Sec
}
