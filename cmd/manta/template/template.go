package template

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "template",
		Short:   "manage templates",
		Aliases: []string{"tmpl"},
	}

	cmd.AddCommand(apply())

	return cmd
}
