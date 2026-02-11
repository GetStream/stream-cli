package feeds

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeds",
		Short: "Allows you to interact with your Feeds applications",
	}

	cmd.AddCommand(NewCmds()...)

	return cmd
}
