package main

import (
	"os"

	"github.com/f1shl3gs/manta/cmd/mantad/launch"
)

func main() {
	rootCmd := launch.Command()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
