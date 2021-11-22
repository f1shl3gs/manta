package bolt

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bolt",
		Short: "analysis or bench the boltdb",
	}

	cmd.AddCommand(bench())

	return cmd
}
