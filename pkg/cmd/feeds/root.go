package feeds

import (
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/cmd/feeds/activity"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/bookmark"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/collection"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/comment"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/feed"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/feedgroup"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/feedview"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/follow"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/membership"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/stats"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/user"
	"github.com/GetStream/stream-cli/pkg/cmd/feeds/visibility"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeds",
		Short: "Allows you to interact with your Feeds applications",
	}

	cmd.AddCommand(activity.NewCmds()...)
	cmd.AddCommand(comment.NewCmds()...)
	cmd.AddCommand(feedgroup.NewCmds()...)
	cmd.AddCommand(feed.NewCmds()...)
	cmd.AddCommand(follow.NewCmds()...)
	cmd.AddCommand(feedview.NewCmds()...)
	cmd.AddCommand(collection.NewCmds()...)
	cmd.AddCommand(bookmark.NewCmds()...)
	cmd.AddCommand(membership.NewCmds()...)
	cmd.AddCommand(visibility.NewCmds()...)
	cmd.AddCommand(stats.NewCmds()...)
	cmd.AddCommand(user.NewCmds()...)

	return cmd
}
