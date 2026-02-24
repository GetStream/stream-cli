package video

import (
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/cmd/video/broadcast"
	"github.com/GetStream/stream-cli/pkg/cmd/video/call"
	"github.com/GetStream/stream-cli/pkg/cmd/video/calltype"
	"github.com/GetStream/stream-cli/pkg/cmd/video/device"
	"github.com/GetStream/stream-cli/pkg/cmd/video/edge"
	"github.com/GetStream/stream-cli/pkg/cmd/video/feedback"
	"github.com/GetStream/stream-cli/pkg/cmd/video/guest"
	"github.com/GetStream/stream-cli/pkg/cmd/video/recording"
	"github.com/GetStream/stream-cli/pkg/cmd/video/sip"
	"github.com/GetStream/stream-cli/pkg/cmd/video/stats"
	"github.com/GetStream/stream-cli/pkg/cmd/video/storage"
	"github.com/GetStream/stream-cli/pkg/cmd/video/transcription"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "video",
		Short: "Allows you to interact with your Video applications",
	}

	cmd.AddCommand(call.NewCmds()...)
	cmd.AddCommand(calltype.NewCmds()...)
	cmd.AddCommand(recording.NewCmds()...)
	cmd.AddCommand(transcription.NewCmds()...)
	cmd.AddCommand(broadcast.NewCmds()...)
	cmd.AddCommand(stats.NewCmds()...)
	cmd.AddCommand(feedback.NewCmds()...)
	cmd.AddCommand(sip.NewCmds()...)
	cmd.AddCommand(device.NewCmds()...)
	cmd.AddCommand(storage.NewCmds()...)
	cmd.AddCommand(guest.NewCmds()...)
	cmd.AddCommand(edge.NewCmds()...)

	return cmd
}
