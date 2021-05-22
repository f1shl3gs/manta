package completion

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate bash/fish/zsh completion scripts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rootCmd := cmd.Parent()
			writer := os.Stdout

			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletion(writer)
			case "fish":
				return rootCmd.GenFishCompletion(writer, true)
			case "zsh":
				return rootCmd.GenZshCompletion(writer)
			default:
				return errors.Errorf("unknown or unsupported shell %q", args[0])
			}
		},
	}

	return cmd
}
