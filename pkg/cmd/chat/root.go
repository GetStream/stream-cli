package chat

import (
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/cmd/chat/app"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/blocklist"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/campaign"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/channel"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/channeltype"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/customcmd"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/device"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/draft"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/events"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/export"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/file"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/imports"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/location"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/message"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/moderation"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/og"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/permission"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/poll"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/push"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/reaction"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/reminder"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/role"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/search"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/segment"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/storage"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/thread"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/unread"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/user"
	"github.com/GetStream/stream-cli/pkg/cmd/chat/watch"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Allows you to interact with your Chat applications",
	}

	cmd.AddCommand(app.NewCmds()...)
	cmd.AddCommand(blocklist.NewCmds()...)
	cmd.AddCommand(campaign.NewCmds()...)
	cmd.AddCommand(channel.NewCmds()...)
	cmd.AddCommand(channeltype.NewCmds()...)
	cmd.AddCommand(customcmd.NewCmds()...)
	cmd.AddCommand(device.NewCmds()...)
	cmd.AddCommand(draft.NewCmds()...)
	cmd.AddCommand(events.NewCmds()...)
	cmd.AddCommand(export.NewCmds()...)
	cmd.AddCommand(file.NewCmds()...)
	cmd.AddCommand(imports.NewCmds()...)
	cmd.AddCommand(location.NewCmds()...)
	cmd.AddCommand(message.NewCmds()...)
	cmd.AddCommand(moderation.NewCmds()...)
	cmd.AddCommand(og.NewCmds()...)
	cmd.AddCommand(permission.NewCmds()...)
	cmd.AddCommand(poll.NewCmds()...)
	cmd.AddCommand(push.NewCmds()...)
	cmd.AddCommand(reaction.NewCmds()...)
	cmd.AddCommand(reminder.NewCmds()...)
	cmd.AddCommand(role.NewCmds()...)
	cmd.AddCommand(search.NewCmds()...)
	cmd.AddCommand(segment.NewCmds()...)
	cmd.AddCommand(storage.NewCmds()...)
	cmd.AddCommand(thread.NewCmds()...)
	cmd.AddCommand(unread.NewCmds()...)
	cmd.AddCommand(user.NewCmds()...)
	cmd.AddCommand(watch.NewCmds()...)

	return cmd
}
