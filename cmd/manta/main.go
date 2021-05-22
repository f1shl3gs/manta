package main

import (
	"os"

	"github.com/f1shl3gs/manta/cmd/manta/completion"
	"github.com/f1shl3gs/manta/cmd/manta/template"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "manta",
		Short:        "manage manta service",
		SilenceUsage: true,
	}

	rootCmd.AddCommand(
		template.New(),
		completion.New(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
