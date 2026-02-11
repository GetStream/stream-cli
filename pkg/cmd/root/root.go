package root

import (
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/cmd/chat"
	cfgCmd "github.com/GetStream/stream-cli/pkg/cmd/config"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/version"
)

var cfgPath = new(string)

func NewCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "stream-cli <command> <subcommand> [flags]",
		Short: "Stream CLI",
		Long:  "Interact with your Stream applications easily",
		Example: heredoc.Doc(`
			# Get Chat application settings
			$ stream-cli chat get-app

			# List all Chat channel types
			$ stream-cli chat list-channel-types

			# Create a new Chat user
			$ stream-cli chat upsert-user --properties "{\"id\":\"my-user-1\"}"
		`),
		Version: version.FmtVersion(),
	}

	fl := root.PersistentFlags()
	fl.String("app", "", "[optional] Application name to use as it's defined in the configuration file")
	fl.StringVar(cfgPath, "config", "", "[optional] Explicit config file path")

	root.AddCommand(
		cfgCmd.NewRootCmd(),
		chat.NewRootCmd(),
		feeds.NewRootCmd(),
	)

	cobra.OnInitialize(config.GetInitConfig(root, cfgPath))

	root.SetOut(os.Stdout)

	return root
}
