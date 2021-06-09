package cluster

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "list, add or remove members",
	}

	cmd.AddCommand(list())
	cmd.AddCommand(add())
	cmd.AddCommand(remove())

	cmd.PersistentFlags().String("host", "127.0.0.1:8088", "")

	return cmd
}

func list() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list all members in cluster",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func add() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add new node to cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}

func remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Short:   "remove member in cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
