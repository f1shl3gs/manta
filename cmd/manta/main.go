package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "manta",
		Short:        "manage manta service",
		SilenceUsage: true,
	}

	// todo: implement

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
