package machine

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	procDir string
)

func GetTotalMemory() (uint64, error) {
	f, err := os.OpenFile(filepath.Join(procDir, "meminfo"), os.O_RDONLY, 0644)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}

		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}

		if fields[0] == "MemTotal" {
			vu := strings.Fields(fields[1])
			val, err := strconv.ParseUint(vu[0], 10, 64)
			if err != nil {
				return 0, err
			}

			return val * 1024, nil
		}
	}

	panic("cannot parse MemTotal")
}
