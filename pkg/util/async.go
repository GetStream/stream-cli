package util

import (
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/urfave/cli/v2"
)

func WaitForAsyncCompletion(ctx *cli.Context, c *stream.Client, taskId string, timeoutSeconds int) error {
	PrintMessage(ctx, "Task ID: "+taskId)
	PrintMessage(ctx, "Waiting for async operation to complete... ⏳")

	for i := 0; i < timeoutSeconds; i++ {
		resp, err := c.GetTask(ctx.Context, taskId)
		if err != nil {
			return err
		}

		if resp.Status == stream.TaskStatusCompleted {
			PrintHappyMessage(ctx, "Async operation completed successfully.")
			return nil
		}

		if i%5 == 0 {
			PrintMessage(ctx, "Still loading... ⏳")
		}

		time.Sleep(time.Second * 1)
	}

	PrintMessage(ctx, "Async operation did not finish it in 10 seconds. "+
		"You'll need to continue polling with `stream-cli watch -t "+taskId+"` command.")

	return nil
}
