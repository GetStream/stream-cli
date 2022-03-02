package chat

import (
	"github.com/GetStream/stream-cli/pkg/cmd/chat/channel"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/watch"
	"github.com/spf13/cobra"
)

func NewChatRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Interact with your Stream Chat application",
	}

	cmd.AddCommand(watch.NewWatchCmd())
	cmd.AddCommand(channel.NewChannelCmds()...)

	return cmd
}
