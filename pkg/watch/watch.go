package watch

import (
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/util"
	"github.com/urfave/cli/v2"
)

func NewWatchCmd(config *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "watch",
		Usage:       "Waits for async operations to complete.",
		Description: "Waits for async operations to complete.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "task-id",
				Aliases:  []string{"t"},
				Usage:    "Task ID to wait for",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "timeout",
				Usage:    "Number of seconds to wait for completion.",
				Value:    30,
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			c, err := config.GetStreamClient(ctx)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			return util.WaitForAsyncCompletion(ctx, c, ctx.String("task-id"), ctx.Int("timeout"))
		},
	}
}
