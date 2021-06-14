package main

import (
	"os"

	"github.com/f1shl3gs/manta/cmd/mantad/cluster"
	"github.com/f1shl3gs/manta/cmd/mantad/completion"
	"github.com/f1shl3gs/manta/cmd/mantad/launcher"
	"github.com/f1shl3gs/manta/cmd/mantad/version"
)

func main() {
	rootCmd := launcher.NewCommand()

	rootCmd.AddCommand(version.New())
	rootCmd.AddCommand(cluster.New())
	rootCmd.AddCommand(completion.New())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
