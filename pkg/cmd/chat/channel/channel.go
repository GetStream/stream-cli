package channel

import (
	"encoding/json"
	"errors"
	"text/tabwriter"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

func NewChannelCmds() []*cobra.Command {
	return []*cobra.Command{
		getChannelCmd(),
		createChannelCmd(),
		deleteChannelCmd(),
		updateChannelCmd(),
		listChannelsCmd()}
}

func getChannelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-channel --type [channel-type] --id [channel-id]",
		Short: "Returns a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			chanId, _ := cmd.Flags().GetString("id")

			r, err := c.Channel(chanType, chanId).Query(cmd.Context(), &stream.QueryRequest{})
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			t := tabby.NewCustom(w)
			t.AddHeader("CID", "Member count", "Created By", "Last Message At", "Created At", "Updated At", "Custom Data")
			t.AddLine(r.Channel.CID,
				r.Channel.MemberCount,
				r.Channel.CreatedBy.ID,
				r.Channel.LastMessageAt.Format(time.RFC822),
				r.Channel.CreatedAt.Format(time.RFC822),
				r.Channel.UpdatedAt.Format(time.RFC822),
				r.Channel.ExtraData)
			t.Print()

			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("id")

	return cmd
}

func createChannelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-channel --type [channel-type] --id [channel-id]",
		Short: "Creates a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			chanId, _ := cmd.Flags().GetString("id")
			user, _ := cmd.Flags().GetString("user")

			r, err := c.CreateChannel(cmd.Context(), chanType, chanId, user, nil)
			if err != nil {
				return err
			}

			if time.Now().UTC().Unix()-r.Channel.CreatedAt.Unix() > 3 {
				return errors.New("channel already exists")
			}

			cmd.Printf("Successfully created channel [%s]", r.Channel.CID)

			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("user", "u", "", "[required] User id")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("user")

	return cmd
}

func deleteChannelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-channel --type [channel-type] --id [channel-id]",
		Short: "Deletes a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			chanId, _ := cmd.Flags().GetString("id")
			hard, _ := cmd.Flags().GetBool("hard")
			cids := []string{chanType + ":" + chanId}

			resp, err := c.DeleteChannels(cmd.Context(), cids, hard)
			if err != nil {
				return err
			}

			if resp.TaskID != "" {
				cmd.Printf("Successfully initiated channel deletion. Task id: %s", resp.TaskID)
			} else {
				cmd.PrintErr("Channel deletion failed")
			}
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.Bool("hard", false, "[optional] Channel will be hard deleted. This action is irrevocable.")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("id")

	return cmd
}

func updateChannelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-channel --type [channel-type] --id [channel-id]",
		Short: "Updates a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			chanId, _ := cmd.Flags().GetString("id")
			p, _ := cmd.Flags().GetString("properties")

			props := make(map[string]interface{})
			err = json.Unmarshal([]byte(p), &props)
			if err != nil {
				return err
			}

			ch := c.Channel(chanType, chanId)
			_, err = ch.Update(cmd.Context(), props, nil)

			cmd.Printf("Successfully updated channel [%s]", chanId)

			return err
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("properties", "p", "", "[required] Channel properties to update")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("properties")

	return cmd
}

func listChannelsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-channels --type [channel-type]",
		Short: "Lists channels",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			offset, _ := cmd.Flags().GetInt("offset")
			limit, _ := cmd.Flags().GetInt("limit")

			resp, err := c.QueryChannels(cmd.Context(), &stream.QueryOption{
				Filter: map[string]interface{}{
					"type": chanType,
				},
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				return err
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			t := tabby.NewCustom(w)
			t.AddHeader("CID", "Member count", "Created By", "Created At", "Updated At")
			for _, c := range resp.Channels {
				t.AddLine(c.CID, c.MemberCount, c.CreatedBy.ID, c.CreatedAt.Format(time.RFC822), c.UpdatedAt.Format(time.RFC822))
			}
			t.Print()

			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.IntP("offset", "o", 0, "[optional] Number of channels to skip during pagination")
	fl.IntP("limit", "l", 10, "[optional] Number of channels to return. Used for pagination")
	cmd.MarkFlagRequired("type")

	return cmd
}
