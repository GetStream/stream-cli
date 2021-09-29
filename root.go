package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "stream-cli [options] <command> <subcommand>",
		Long: `Manage seamlessly Stream applications from the command line`,
		//SilenceErrors: true,
		//SilenceUsage: true,
		Run: func(_ *cobra.Command, _ []string) {},
	}

	cmd.AddCommand(NewRootConfigCmd())

	return cmd
}

