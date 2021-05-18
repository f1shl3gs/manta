package wal

import (
	"os"
	"sync"
)

type WAL struct {
	mtx  sync.Mutex
	file *os.File
}
