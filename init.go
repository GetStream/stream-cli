package cli

import "github.com/urfave/cli/v2"

func NewInitCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "init",
		Usage:       "Initializes application credentials.",
		UsageText:   "stream-cli init",
		Description: "Initializes application credentials.",

		Action: func(ctx *cli.Context) error {
			return RunQuestionnaire(config)
		},
	}
}
