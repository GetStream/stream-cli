package root

import (
	"fmt"
	"time"

	"github.com/GetStream/stream-cli/pkg/channel"
	cfg "github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/update"
	"github.com/GetStream/stream-cli/pkg/util"
	"github.com/GetStream/stream-cli/pkg/version"
	"github.com/GetStream/stream-cli/pkg/watch"
	"github.com/urfave/cli/v2"
)

func NewRootCmd(config *cfg.Config) *cli.App {
	return &cli.App{
		Name:        "stream-cli",
		Usage:       "Interact with your Stream applications easily",
		Description: "The official Stream CLI allows you to interact with your applications easily",
		Compiled:    time.Now(),
		Copyright:   fmt.Sprintf("(c) %d Stream.io Inc.", time.Now().Year()),
		Version:     version.FmtVersion(),
		ExitErrHandler: func(c *cli.Context, err error) {
			if exitErr, ok := err.(cli.ExitCoder); ok {
				if err.Error() != "" {
					util.PrintSadMessage(c, err.Error())
				}
				cli.OsExiter(exitErr.ExitCode())
			}
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "app",
				Usage: "The name of the application to use. If not defined the default application will be used.",
			},
		},
		Commands: []*cli.Command{
			cfg.NewRootConfigCmd(config),
			channel.NewChannelCmd(config),
			update.NewUpdateCmd(config),
			watch.NewWatchCmd(config),
		},
	}
}
