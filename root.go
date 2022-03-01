package cli

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

func NewRootCmd(config *Config) *cli.App {
	return &cli.App{
		Name:        "stream-cli",
		Usage:       "Interact with your Stream applications easily",
		Description: "The official Stream CLI allows you to interact with your applications easily",
		Compiled:    time.Now(),
		Copyright:   fmt.Sprintf("(c) %d Stream.io Inc.", time.Now().Year()),
		Version:     fmtVersion(),
		ExitErrHandler: func(c *cli.Context, err error) {
			if exitErr, ok := err.(cli.ExitCoder); ok {
				if err.Error() != "" {
					PrintSadMessage(c, err.Error())
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
			NewRootConfigCmd(config),
			NewInitCmd(config),
			NewChannelCmd(config),
			NewUpdateCmd(config),
			NewWatchCmd(config),
		},
	}
}
