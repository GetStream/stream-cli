package channel

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getCmd(),
		createCmd(),
		deleteCmd(),
		updateCmd(),
		updatePartialCmd(),
		listCmd()}
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-channel --type [channel-type] --id [channel-id]",
		Short: "Return a channel",
		Example: heredoc.Doc(`
			# Returns 'redteam' channel of 'messaging' channel type as JSON
			$ stream-cli chat get-channel --type messaging --id redteam

			# Returns 'blueteam' channel of 'messaging' channel type as a browsable tree
			$ stream-cli chat get-channel --type messaging --id blueteam --output-format tree
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			chanId, _ := cmd.Flags().GetString("id")

			r, err := c.Channel(chanType, chanId).Query(cmd.Context(), &stream.QueryRequest{})
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, r.Channel)
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("id")

	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-channel --type [channel-type] --id [channel-id] --user [user-id]",
		Short: "Create a channel",
		Long: heredoc.Doc(`
			This command allows you to create a new channel. If it
			exists already an error will be thrown.
		`),
		Example: heredoc.Doc(`
			# Create a channel with id 'redteam' of type 'messaging' by 'joe'
			$ stream-cli chat create-channel --type messaging --id redteam --user joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
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

			cmd.Printf("Successfully created channel [%s]\n", r.Channel.CID)

			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("user", "u", "", "[required] User id who will be considered as the creator of the channel")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("user")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-channel --type [channel-type] --id [channel-id]",
		Short: "Delete a channel",
		Long: heredoc.Doc(`
			This command allows you to delete a channel. This operation is asynchronous
			in the backend so a task id is returned. You need to use the watch
			commnand to poll the results.
		`),
		Example: heredoc.Doc(`
			# Delete a channel with id 'redteam' of type 'messaging'
			$ stream-cli chat delete-channel --type messaging --id redteam
			> Successfully initiated channel deletion. Task id: 66bbcdcd-b133-43ce-ab63-557c14d2a168

			# Wait for the task to complete
			$ stream-cli chat watch 66bbcdcd-b133-43ce-ab63-557c14d2a168
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
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
				cmd.Printf("Successfully initiated channel deletion. Task id: %s\n", resp.TaskID)
			} else {
				return errors.New("channel deletion failed")
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

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-channel --type [channel-type] --id [channel-id] --properties [raw-json-properties]",
		Short: "Update a channel",
		Long: heredoc.Doc(`
			Updates an existing channel. The 'properties' are specified as a raw json string. The valid
			properties are the 'ChannelRequest' object of the official documentation.
			Such as 'team', 'frozen', 'disabled' or any custom property.
			https://getstream.io/chat/docs/rest/#channels-updatechannel
		`),
		Example: heredoc.Doc(`
			# Unfreeze a channel
			$ stream-cli chat update-channel --type messaging --id redteam --properties "{\"frozen\":false}"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
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
			if err != nil {
				return err
			}

			cmd.Printf("Successfully updated channel [%s]\n", chanId)
			return nil
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

func updatePartialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-channel-partial --type [channel-type] --id [channel-id] --set [key-value-pairs] --unset [property-names]",
		Short: "Update a channel partially",
		Long: heredoc.Doc(`
			Updates an existing channel. The 'set' property is a comma separated list of key value pairs.
			The 'unset' property is a comma separated list of property names.
		`),
		Example: heredoc.Doc(`
			# Freeze a channel and set 'age' to 21. At the same time, remove 'haircolor' and 'height'.
			stream-cli chat update-channel-partial --type messaging --id channel1 --set frozen=true,age=21 --unset haircolor,height
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			chanId, _ := cmd.Flags().GetString("id")
			set, _ := cmd.Flags().GetStringToString("set")
			unset, _ := cmd.Flags().GetString("unset")

			s := make(map[string]interface{}, len(set))
			for k, v := range set {
				s[k] = v
			}

			u := make([]string, 0)
			for _, v := range strings.Split(unset, ",") {
				if v != "" {
					u = append(u, strings.TrimSpace(v))
				}
			}

			ch := c.Channel(chanType, chanId)
			_, err = ch.PartialUpdate(cmd.Context(), stream.PartialUpdate{Set: s, Unset: u})
			if err != nil {
				return err
			}

			cmd.Printf("Successfully updated channel [%s]\n", chanId)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringToStringP("set", "s", map[string]string{}, "[optional] Comma-separated key-value pairs to set")
	fl.StringP("unset", "u", "", "[optional] Comma separated list of properties to unset")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("id")

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-channels --type [channel-type]",
		Short: "List channels",
		Long: heredoc.Doc(`
			List all channels of a given channel type. You can also provide
			a limit for paginating the results.
		`),
		Example: heredoc.Doc(`
			# List the top 5 'messaging' channels as a json
			$ stream-cli chat list-channels --type messaging --limit 5

			# List the top 20 'livestream' channels as a browsable tree
			$ stream-cli chat list-channels --type livestream --limit 20 --output-format tree
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			limit, _ := cmd.Flags().GetInt("limit")

			resp, err := c.QueryChannels(cmd.Context(), &stream.QueryOption{
				Filter: map[string]interface{}{
					"type": chanType,
				},
				Sort:  []*stream.SortOption{{Field: "cid", Direction: 1}},
				Limit: limit,
			})
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp)
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.IntP("limit", "l", 10, "[optional] Number of channels to return. Used for pagination")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	cmd.MarkFlagRequired("type")

	return cmd
}
