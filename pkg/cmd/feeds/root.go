package feeds

import (
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/cmd/feeds/imports"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeds",
		Short: "Allows you to interact with your Feeds applications",
	}

	cmd.AddCommand(imports.NewCmds()...)

	return cmd
}
