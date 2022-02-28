package cli

import (
	"encoding/json"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/cheynewallace/tabby"
	"github.com/urfave/cli/v2"
)

func NewChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "channel",
		Usage:       "Manage channels.",
		Description: "Manage channels.",
		Subcommands: []*cli.Command{
			getChannelCmd(config),
			newListChannelsCmd(config),
			newCreateChannelCmd(config),
			newDeleteChannelCmd(config),
			newUpdateChannelCmd(config),
		},
	}
}

func getChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "get",
		Usage:       "Get a channel by channel type and channel name.",
		Description: "Get a channel by channel type and channel name.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "channel-type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "channel-name",
				Aliases:  []string{"n"},
				Usage:    "The name of the channel.",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "json",
				Usage:    "Print the raw JSON representation of the API response.",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			r, err := getOrCreateChannel(config, ctx)

			if err != nil {
				return err
			}

			t := tabby.New()
			t.AddHeader("CID", "Member count", "Created By", "Created At", "Updated At", "Extra Data")
			t.AddLine(r.Channel.CID,
				r.Channel.MemberCount,
				r.Channel.CreatedBy.ID,
				r.Channel.CreatedAt.Format(time.RFC822),
				r.Channel.UpdatedAt.Format(time.RFC822),
				r.Channel.ExtraData)

			if ctx.Bool("json") {
				PrintRawJson(r)
			} else {
				t.Print()
			}

			return nil
		},
	}
}

func newCreateChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create a new channel or return an existing one.",
		UsageText: "stream-cli channel create --channel-type [channel-type] --channel-name [channel-name]",
		Description: "Creates a new channel or returns an existing one if it already exists. " +
			"If channel name is not provided it will be generated automatically.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "channel-type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "channel-name",
				Aliases:  []string{"n"},
				Usage:    "The name of the channel.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "The user ID of the user who will be the creator of the channel.",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "json",
				Usage:    "Print the raw JSON representation of the API response.",
				Required: false,
			},
		},

		Action: func(ctx *cli.Context) error {
			r, err := getOrCreateChannel(config, ctx)

			if err != nil {
				return err
			}

			PrintHappyMessageFormatted("Successfully created channel [%s].", r.Channel.CID)

			if ctx.Bool("json") {
				PrintRawJson(r)
			}

			return nil
		},
	}
}

func getOrCreateChannel(config *Config, ctx *cli.Context) (*stream.CreateChannelResponse, error) {
	t := ctx.String("channel-type")
	n := ctx.String("channel-name")
	u := ctx.String("user")

	c, err := config.GetStreamClient(ctx)

	if err != nil {
		return nil, cli.Exit(err.Error(), 1)
	}

	if u == "" {
		// This means that it's a GetChannel operation not a CreateChannel one. Let's use a dummy name then.
		u = "stream-go-cli"
	}

	r, err := c.CreateChannel(ctx.Context, t, n, u, nil)

	if err != nil {
		return nil, cli.Exit(err.Error(), 1)
	}

	return r, nil
}

func newDeleteChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Usage:       "Delete a channel.",
		UsageText:   "stream-cli channel delete --channel-type messaging --channel-name my-team-channel --hard",
		Description: "Delete a channel.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "channel-type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "channel-name",
				Aliases:  []string{"n"},
				Usage:    "The name of the channel.",
				Required: true,
			},
			&cli.BoolFlag{
				Name: "hard",
				// "h" alias is already used by "--help" flag, so we don't have an alias here
				Usage:    "Hard deleted channels cannot be restored.",
				Value:    false,
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "json",
				Usage:    "Print the raw JSON representation of the API response.",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			c, err := config.GetStreamClient(ctx)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			hard := ctx.Bool("hard")
			cids := []string{ctx.String("channel-type") + ":" + ctx.String("channel-name")}

			resp, err := c.DeleteChannels(ctx.Context, cids, hard)

			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if ctx.Bool("json") {
				PrintRawJson(resp)
			}

			if resp.TaskID == "" {
				PrintMessage("The channel is already deleted.")
				return nil
			}

			return WaitForAsyncCompletion(ctx.Context, c, resp.TaskID)
		},
	}
}

func newUpdateChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "update",
		Usage:       "Update a channel.",
		UsageText:   "stream-cli channel update --channel-type messaging --channel-name my-team-channel --properties '{\"frozen\": \"true\"}'",
		Description: "Update a channel.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "channel-type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "channel-name",
				Aliases:  []string{"n"},
				Usage:    "The name of the channel.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "properties",
				Aliases:  []string{"p"},
				Usage:    "Properties of the update, as a raw JSON string. Example: '{\"frozen\": \"true\"}'.",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "json",
				Usage:    "Print the raw JSON representation of the API response.",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			c, err := config.GetStreamClient(ctx)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			t := ctx.String("channel-type")
			n := ctx.String("channel-name")
			p := ctx.String("properties")
			props := make(map[string]interface{})
			json.Unmarshal([]byte(p), &props)

			ch := c.Channel(t, n)
			_, err = ch.Update(ctx.Context, props, nil)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			PrintHappyMessageFormatted("Successfully updated channel [%s].", n)

			if ctx.Bool("json") {
				PrintRawJson(ch)
			}

			return nil
		},
	}
}

func newListChannelsCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all channels.",
		UsageText:   "stream-cli channel list --channel-type messaging",
		Description: "Update a channel.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "channel-type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'. If not specified, it returns all channels.",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "limit",
				Aliases:  []string{"l"},
				Usage:    "The number of channels to return. Use for pagination.",
				Value:    50,
				Required: false,
			},
			&cli.IntFlag{
				Name:     "offset",
				Aliases:  []string{"o"},
				Usage:    "The number of channels to skip during pagination.",
				Value:    0,
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "json",
				Usage:    "Print the raw JSON representation of the API response.",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			c, err := config.GetStreamClient(ctx)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			resp, err := c.QueryChannels(ctx.Context, &stream.QueryOption{
				Filter: map[string]interface{}{
					"type": ctx.String("channel-type"),
				},
				Limit:  ctx.Int("limit"),
				Offset: ctx.Int("offset"),
			})
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			t := tabby.New()
			t.AddHeader("CID", "Member count", "Created By", "Created At", "Updated At")
			for _, c := range resp.Channels {
				t.AddLine(c.CID, c.MemberCount, c.CreatedBy.ID, c.CreatedAt.Format(time.RFC822), c.UpdatedAt.Format(time.RFC822))
			}

			if ctx.Bool("json") {
				PrintRawJson(resp)
			} else {
				t.Print()
			}

			return nil
		},
	}
}
