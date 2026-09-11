package gopassxc

import (
	"fmt"
	"os"

	helpers "github.com/iFleksy/gopassxc/cmd/helpers"
	"github.com/spf13/cobra"
)

type GetLoginsCMDOpts struct {
	TOTP bool
}

func (o *GetLoginsCMDOpts) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.TOTP, "totp", false, "Get TOTP code")
}

func GetLoginCMD() *cobra.Command {
	opts := GetLoginsCMDOpts{}
	cmd := &cobra.Command{
		Use: "login <url>",
		Example: `login search www.yandex.ru
login search www.yandex.ru --totp
		`,
		Short:        "Manage logins",
		Long:         `Manage logins stored in KeePassXC.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := helpers.BuildWithConfig()
			if err != nil {
				return err
			}

			client, err = helpers.PrepareClient(cmd.Context(), client)
			if err != nil {
				return err
			}

			err = client.TestAssociate(cmd.Context(), true)
			if err != nil {
				return err
			}

			// Use the client here
			clients, err := client.GetLogins(cmd.Context(), args[0])

			if err != nil {
				return err
			}

			if len(clients) > 1 {
				return fmt.Errorf("found many records")
			}

			record := clients[0]

			if opts.TOTP {
				fmt.Fprintln(os.Stdout, record.Totp)
			} else {
				fmt.Fprintln(os.Stdout, record.Password)
			}
			return nil
		},
	}
	opts.AddFlags(cmd)
	return cmd
}
