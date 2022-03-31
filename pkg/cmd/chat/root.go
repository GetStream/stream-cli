package chat

import (
	"github.com/GetStream/stream-cli/pkg/cmd/app"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/channel"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/channeltype"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/file"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/imports"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/message"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/user"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/watch"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Allows you to interact with your Chat applications",
	}

	cmd.AddCommand(app.NewCmds()...)
	cmd.AddCommand(channel.NewCmds()...)
	cmd.AddCommand(channeltype.NewCmds()...)
	cmd.AddCommand(file.NewCmds()...)
	cmd.AddCommand(imports.NewCmds()...)
	cmd.AddCommand(message.NewCmds()...)
	cmd.AddCommand(user.NewCmds()...)
	cmd.AddCommand(watch.NewCmds()...)

	return cmd
}
