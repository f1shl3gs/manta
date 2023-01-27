package main

import (
	"os"

	"github.com/f1shl3gs/manta/cmd/mantad/launch"
	"github.com/f1shl3gs/manta/cmd/mantad/version"
)

func main() {
	rootCmd := launch.Command()

	rootCmd.AddCommand(version.Command())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
