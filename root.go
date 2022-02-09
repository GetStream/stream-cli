package cli

import (
	"github.com/urfave/cli/v2"
)

func NewRootCmd() *cli.App {
	app := &cli.App{
		Name:        "stream-cli",
		Usage:       "Interact with your Stream applications easily",
		Description: "The official Stream CLI allows you to interact with your applications easily",
		Commands: []*cli.Command{
			NewRootConfigCmd(),
		},
	}

	return app
}
