package video

import (
	"github.com/spf13/cobra"

	rawrecording "github.com/GetStream/stream-cli/pkg/cmd/raw-recording-tool"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "video",
		Short: "Video processing commands",
	}
	cmd.AddCommand(rawrecording.NewRootCmd())
	return cmd
}
