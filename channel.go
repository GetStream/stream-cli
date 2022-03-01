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
			listChannelsCmd(config),
			createChannelCmd(config),
			deleteChannelCmd(config),
			updateChannelCmd(config),
		},
	}
}

func getChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "get",
		Usage:       "Get a channel by channel type and channel name.",
		UsageText:   "stream-cli channel get --type [channel-type] --id [channel-id]",
		Description: "Get a channel by channel type and channel name.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "The id of the channel.",
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

			r, err := c.Channel(ctx.String("type"), ctx.String("id")).Query(ctx.Context, &stream.QueryRequest{})
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if ctx.Bool("json") {
				PrintRawJson(ctx, r)
			} else {
				t := tabby.New()
				t.AddHeader("CID", "Member count", "Created By", "Created At", "Updated At", "Custom Data")
				t.AddLine(r.Channel.CID,
					r.Channel.MemberCount,
					r.Channel.CreatedBy.ID,
					r.Channel.CreatedAt.Format(time.RFC822),
					r.Channel.UpdatedAt.Format(time.RFC822),
					r.Channel.ExtraData)
				t.Print()
			}

			return nil
		},
	}
}

func createChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "create",
		Usage:       "Create a new channel or return an existing one.",
		UsageText:   "stream-cli channel create --type [channel-type] --id [channel-id]",
		Description: "Creates a new channel or returns an existing one if it already exists.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "The id of the channel.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "The ID of the user who will be the creator of the channel.",
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

			r, err := c.CreateChannel(ctx.Context, ctx.String("type"), ctx.String("id"), ctx.String("user"), nil)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if time.Now().UTC().Unix()-r.Channel.CreatedAt.Unix() > 3 {
				return cli.Exit("The channel exists already", 1)
			}

			PrintHappyMessageFormatted(ctx, "Successfully created channel [%s].", r.Channel.CID)

			if ctx.Bool("json") {
				PrintRawJson(ctx, r)
			}

			return nil
		},
	}
}

func deleteChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Usage:       "Delete a channel.",
		UsageText:   "stream-cli channel delete --type [channel-type] --id [channel-id] --hard",
		Description: "Delete a channel.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "The id of the channel.",
				Required: true,
			},
			&cli.BoolFlag{
				Name: "hard",
				// "h" alias is already used by "--help" flag, so we don't have an alias here
				Usage:    "Channel will be hard deleted. This action is irrevocable.",
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
			cids := []string{ctx.String("type") + ":" + ctx.String("id")}

			resp, err := c.DeleteChannels(ctx.Context, cids, hard)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			if ctx.Bool("json") {
				PrintRawJson(ctx, resp)
			}

			if resp.TaskID == "" {
				PrintMessage(ctx, "The channel is already deleted.")
				return nil
			}

			return WaitForAsyncCompletion(ctx, c, resp.TaskID, 10)
		},
	}
}

func updateChannelCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "update",
		Usage:       "Update a channel.",
		UsageText:   "stream-cli channel update --type [channel-type] --id [channel-id] --properties '{\"frozen\": \"true\"}'",
		Description: "Update a channel.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "The type of the channel. Such as 'messaging' or 'livestream'.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "The id of the channel.",
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

			t := ctx.String("type")
			n := ctx.String("id")
			p := ctx.String("properties")
			props := make(map[string]interface{})
			err = json.Unmarshal([]byte(p), &props)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			ch := c.Channel(t, n)
			_, err = ch.Update(ctx.Context, props, nil)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			PrintHappyMessageFormatted(ctx, "Successfully updated channel [%s].", n)

			if ctx.Bool("json") {
				PrintRawJson(ctx, ch)
			}

			return nil
		},
	}
}

func listChannelsCmd(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List all channels.",
		UsageText:   "stream-cli channel list --type [channel-type]",
		Description: "List all channels.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "type",
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
					"type": ctx.String("type"),
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
				PrintRawJson(ctx, resp)
			} else {
				t.Print()
			}

			return nil
		},
	}
}
