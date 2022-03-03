package utils

import (
	"errors"
	"fmt"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/spf13/cobra"
)

func WaitForAsyncCompletion(cmd *cobra.Command, c *stream.Client, taskID string, timeoutSeconds int) error {
	cmd.Println(fmt.Sprintf("Task id: %s\n", taskID))
	cmd.Println("Waiting for async task to complete...⏳")

	for i := 0; i < timeoutSeconds; i++ {
		resp, err := c.GetTask(cmd.Context(), taskID)
		if err != nil {
			return err
		}

		if resp.Status == stream.TaskStatusCompleted {
			cmd.Print("Async operation completed successfully")
			return nil
		}
		if resp.Status == stream.TaskStatusFailed {
			return errors.New("async operation failed")
		}

		if i%5 == 0 {
			cmd.Print("Still loading... ⏳")
		}

		time.Sleep(time.Second)
	}

	cmd.PrintErrf("Async operation timed out after [%d] seconds", timeoutSeconds)

	return nil
}
