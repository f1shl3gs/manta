package version

import (
	"fmt"

	"github.com/f1shl3gs/manta/version"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("version:     ", version.Version)
			fmt.Println("commit:      ", version.GitSHA)
			fmt.Println("branch:      ", version.GitBranch)
			return nil
		},
	}
}
