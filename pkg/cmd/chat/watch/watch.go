package watch

import (
	"github.com/GetStream/stream-cli/pkg/cmd/chat/utils"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch [task-id]",
		Short: "Wait for an async task to complete",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			timeout, _ := cmd.Flags().GetInt("timeout")
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			return utils.WaitForAsyncCompletion(cmd, c, taskID, timeout)
		},
	}

	cmd.Flags().IntP("timeout", "t", 30, "[optional] Timeout in seconds. Default is 30")

	return []*cobra.Command{cmd}
}
