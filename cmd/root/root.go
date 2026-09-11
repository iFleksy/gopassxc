package root

import (
	"github.com/iFleksy/gopassxc/cmd/gopassxc"
	"github.com/spf13/cobra"
)

func GetRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gopassxc",
		Short: "gopassxc is a command-line tool for managing passwords with KeePassXC.",
		Long:  `gopassxc is a command-line tool for managing passwords with KeePassXC.`,
	}

	cmd.AddCommand(gopassxc.GetLoginCMD())
	cmd.AddCommand(gopassxc.GetDBHashCMD())
	return cmd
}
