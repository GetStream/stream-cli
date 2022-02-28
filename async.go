package cli

import (
	"context"
	"errors"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
)

func WaitForAsyncCompletion(ctx context.Context, c *stream.Client, taskId string) error {
	PrintMessage("Waiting for async operation to complete... ⏳")

	for i := 0; i < 20; i++ {
		resp, err := c.GetTask(ctx, taskId)
		if err != nil {
			return err
		}

		if resp.Status == stream.TaskStatusCompleted {
			PrintHappyMessage("Async operation completed successfully.")
			return nil
		}

		if i%5 == 0 {
			PrintMessage("Still loading... ⏳")
		}

		time.Sleep(time.Millisecond * 500)
	}

	return errors.New("async operation did not complete within 10 seconds")
}
