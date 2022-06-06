package main

import (
	"github.com/f1shl3gs/manta/cmd/mantad/launch"
	"os"
)

func main() {
	rootCmd := launch.Command()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
