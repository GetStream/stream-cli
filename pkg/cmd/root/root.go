package root

import (
	"github.com/GetStream/stream-cli/pkg/cmd/chat"
	cfgCmd "github.com/GetStream/stream-cli/pkg/cmd/config"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/version"
	"github.com/spf13/cobra"
)

var (
	cfgPath *string = new(string)
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "stream-cli <command> <subcommand> [flags]",
		Short:   "Stream CLI",
		Long:    "Interact with your Stream applications easily",
		Version: version.FmtVersion(),
	}

	fl := root.PersistentFlags()
	fl.String("app", "", "[optional] Application id to use as it's defined in the configuration file")
	fl.StringVar(cfgPath, "config", "", "[optional] Explicit config file path")

	root.AddCommand(
		cfgCmd.NewRootConfigCmd(),
		chat.NewChatRootCmd(),
	)

	cobra.OnInitialize(config.GetInitConfig(root, cfgPath))

	return root
}
