package importcmd

import (
	"github.com/spf13/cobra"

	chatImports "github.com/GetStream/stream-cli/pkg/cmd/chat/imports"
	feedsImports "github.com/GetStream/stream-cli/pkg/cmd/feeds/imports"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import data into your Stream applications",
	}

	chatCmd := &cobra.Command{
		Use:   "chat",
		Short: "Import data into Chat",
	}
	chatCmd.AddCommand(chatImports.NewCmds()...)

	feedsCmd := &cobra.Command{
		Use:   "feeds",
		Short: "Import data into Feeds",
	}
	feedsCmd.AddCommand(feedsImports.NewCmds()...)

	cmd.AddCommand(chatCmd, feedsCmd)

	return cmd
}
