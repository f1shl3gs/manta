package version

import (
	"fmt"

	"github.com/f1shl3gs/manta"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version:      ", manta.Version)
			fmt.Println("Branch:       ", manta.Branch)
			fmt.Println("Commit:       ", manta.Commit)
		},
	}
}
