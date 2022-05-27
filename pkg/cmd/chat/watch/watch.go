package watch

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/cmd/chat/utils"
	"github.com/GetStream/stream-cli/pkg/config"
)

func NewCmds() []*cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch [task-id]",
		Short: "Wait for an async task to complete",
		Long: heredoc.Doc(`
			This command waits for a specific async backend operation
			to complete. Such as deleting a user or exporting a channel.
		`),
		Example: heredoc.Doc(`
			# Delete user and watching it complete
			$ stream-cli chat delete-users "my-user-1"
			> Successfully initiated user deletion. Task id: 7586fa0d-dc8d-4f6f-be2d-f952d0e26167

			# Waiting for the task to complete
			$ stream-cli chat watch 7586fa0d-dc8d-4f6f-be2d-f952d0e26167

			# Providing a timeout of 80 seconds
			$ stream-cli chat watch 7586fa0d-dc8d-4f6f-be2d-f952d0e26167 --timeout 80
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			timeout, _ := cmd.Flags().GetInt("timeout")
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			return utils.WaitForAsyncCompletion(cmd, c, taskID, timeout)
		},
	}

	cmd.Flags().IntP("timeout", "t", 30, "[optional] Timeout in seconds. Default is 30")

	return []*cobra.Command{cmd}
}
