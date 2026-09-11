package gopassxc

import (
	"fmt"
	"os"

	helpers "github.com/iFleksy/gopassxc/cmd/helpers"
	"github.com/spf13/cobra"
)

func GetDBHashCMD() *cobra.Command {
	opts := GetLoginsCMDOpts{}
	cmd := &cobra.Command{
		Use:     "db hash",
		Example: `db hash`,
		Short:   "Get database hash",
		Long:    `Get the hash of the KeePassXC database.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := helpers.BuildWithConfig()
			if err != nil {
				return err
			}
			client, err = helpers.PrepareClient(cmd.Context(), client)
			if err != nil {
				return err
			}
			dbHash, err := client.GetDBHash(cmd.Context())
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "%s", dbHash)
			return nil
		},
	}
	opts.AddFlags(cmd)
	return cmd
}
