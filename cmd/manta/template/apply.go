package template

import (
	"fmt"

	"github.com/f1shl3gs/manta"
	"github.com/spf13/cobra"
)

func apply() *cobra.Command {
	var (
		from string
		org  string
	)

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply template from file or url",
		RunE: func(cmd *cobra.Command, args []string) error {
			var orgID manta.ID

			err := orgID.DecodeFromString(org)
			if err != nil {
				return err
			}

			return applyTemplate(orgID, from)
		},
	}

	cmd.Flags().StringVarP(&from, "from", "f", "", "specify the resource uri")
	err := cmd.MarkFlagRequired("from")
	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVarP(&org, "org", "o", "", "specify the organization ID")
	err = cmd.MarkFlagRequired("org")
	if err != nil {
		panic(err)
	}

	return cmd
}

func applyTemplate(orgID manta.ID, from string) error {
	fmt.Println("todo apply template to a specify organization")
	return nil
}
