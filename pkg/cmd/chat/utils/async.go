package utils

import (
	"fmt"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/spf13/cobra"
)

func WaitForAsyncCompletion(cmd *cobra.Command, c *stream.Client, taskId string, timeoutSeconds int) error {
	cmd.Println(fmt.Sprintf("Task id: %s\n", taskId))
	cmd.Println("Waiting for async task to complete...⏳")

	for i := 0; i < timeoutSeconds; i++ {
		resp, err := c.GetTask(cmd.Context(), taskId)
		if err != nil {
			return err
		}

		if resp.Status == stream.TaskStatusCompleted {
			cmd.Print("Async operation completed successfully")
			return nil
		}

		if i%5 == 0 {
			cmd.Print("Still loading... ⏳")
		}

		time.Sleep(time.Second * 1)
	}

	cmd.PrintErrf("Async operation timed out after [%d] seconds", timeoutSeconds)

	return nil
}
