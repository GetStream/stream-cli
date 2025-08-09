package channel

import (
	"encoding/json"
	"errors"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getCmd(),
		createCmd(),
		deleteCmd(),
		updateCmd(),
		updatePartialCmd(),
		listCmd(),
		addMembersCmd(),
		removeMemberCmd(),
		promoteModeratorCmd(),
		demoteModeratorCmd(),
		assignRoleCmd(),
		hideCmd(),
		showCmd(),
		truncateChannelCmd(),
	}
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
			chanID, _ := cmd.Flags().GetString("id")

			r, err := c.Channel(chanType, chanID).Query(cmd.Context(), &stream.QueryRequest{})
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
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-channel --type [channel-type] --id [channel-id] --user [user-id] --properties [raw-json]",
		Short: "Create a channel",
		Long: heredoc.Doc(`
			This command allows you to create a new channel. If it
			exists already an error will be thrown.
		`),
		Example: heredoc.Doc(`
			# Create a channel with id 'redteam' of type 'messaging' by 'joe'
			$ stream-cli chat create-channel --type messaging --id redteam --user joe

			# Create a channel with id 'blueteam' of type 'messaging' by 'joe' with extra data
			$ stream-cli chat create-channel --type messaging --id blueteam --user joe --properties "{\"age\":\"28\"}"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			chanID, _ := cmd.Flags().GetString("id")
			user, _ := cmd.Flags().GetString("user")
			rawProps, _ := cmd.Flags().GetString("properties")

			var props stream.ChannelRequest
			if rawProps != "" {
				err := json.Unmarshal([]byte(rawProps), &props)
				if err != nil {
					return err
				}
			}

			r, err := c.CreateChannel(cmd.Context(), chanType, chanID, user, &props)
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
	fl.StringP("properties", "p", "", "[optional] JSON string of channel properties")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user")

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
			chanID, _ := cmd.Flags().GetString("id")
			hard, _ := cmd.Flags().GetBool("hard")
			cids := []string{chanType + ":" + chanID}

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
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")

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
			chanID, _ := cmd.Flags().GetString("id")
			p, _ := cmd.Flags().GetString("properties")

			props := make(map[string]interface{})
			err = json.Unmarshal([]byte(p), &props)
			if err != nil {
				return err
			}

			ch := c.Channel(chanType, chanID)
			_, err = ch.Update(cmd.Context(), props, nil)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully updated channel [%s]\n", chanID)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("properties", "p", "", "[required] Channel properties to update")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}

func updatePartialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-channel-partial --type [channel-type] --id [channel-id] --set [raw-json] --unset [property-names]",
		Short: "Update a channel partially",
		Long: heredoc.Doc(`
			Updates an existing channel. The 'set' property is a comma separated list of key value pairs.
			The 'unset' property is a comma separated list of property names.
		`),
		Example: heredoc.Doc(`
			# Freeze a channel and set 'age' to 21. At the same time, remove 'haircolor' and 'height'.
			$ stream-cli chat update-channel-partial --type messaging --id channel1 --set '{"frozen":true,"age":21}' --unset haircolor,height
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			flags := cmd.Flags()
			chanType, _ := flags.GetString("type")
			chanID, _ := flags.GetString("id")
			update, err := utils.GetPartialUpdateParam(flags)
			if err != nil {
				return err
			}

			ch := c.Channel(chanType, chanID)
			_, err = ch.PartialUpdate(cmd.Context(), update)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully updated channel [%s]\n", chanID)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("set", "s", "", "[optional] Raw JSON of key-value pairs to set")
	fl.StringP("unset", "u", "", "[optional] Comma separated list of properties to unset")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")

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
	_ = cmd.MarkFlagRequired("type")

	return cmd
}

func addMembersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-members --type [channel-type] --id [channel-id] [user-id-1] [user-id-2] ...",
		Short: "Add members to a channel",
		Example: heredoc.Doc(`
			# Add members joe, jill and jane to 'red-team' channel
			$ stream-cli chat add-members --type messaging --id red-team joe jill jane
		`),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")

			_, err = c.Channel(chType, chID).AddMembers(cmd.Context(), args)
			if err != nil {
				return err
			}

			cmd.Println("Successfully added user(s) to channel")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func removeMemberCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-members --type [channel-type] --id [channel-id] [user-id-1] [user-id-2] ...",
		Short: "Remove members from a channel",
		Example: heredoc.Doc(`
			# Remove members joe, jill and jane from 'red-team' channel
			$ stream-cli chat remove-members --type messaging --id red-team joe jill jane
		`),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")

			_, err = c.Channel(chType, chID).RemoveMembers(cmd.Context(), args, nil)
			if err != nil {
				return err
			}

			cmd.Println("Successfully removed user(s) from channel")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func promoteModeratorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promote-moderators --type [channel-type] --id [channel-id] [user-id-1] [user-id-2] ...",
		Short: "Promote users to channel moderator role",
		Args:  cobra.MinimumNArgs(1),
		Example: heredoc.Doc(`
			# Promote 4 users to moderator
			$ stream-cli chat promote-moderators --type messaging --id red-team joe mike jane jill
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")

			_, err = c.Channel(chType, chID).AddModerators(cmd.Context(), args...)
			if err != nil {
				return err
			}

			cmd.Println("Successfully promoted users to moderators")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func demoteModeratorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demote-moderators --type [channel-type] --id [channel-id] [user-id-1] [user-id-2] ...",
		Short: "Demote users from moderator role",
		Args:  cobra.MinimumNArgs(1),
		Example: heredoc.Doc(`
			# Demote 4 users from moderator role
			$ stream-cli chat demote-moderators --type messaging --id red-team joe mike jane jill
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")

			_, err = c.Channel(chType, chID).DemoteModerators(cmd.Context(), args...)
			if err != nil {
				return err
			}

			cmd.Println("Successfully demoted users from moderators")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}

func assignRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign-role --type [channel-type] --id [channel-id] --user-id [user-id] --role [channel-role-name]",
		Short: "Assign a role to a user",
		Example: heredoc.Doc(`
			# Assign 'channel_moderator' role to user 'joe'
			$ stream-cli chat assign-role --type messaging --id red-team --user-id joe --role channel_moderator
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")
			userID, _ := cmd.Flags().GetString("user-id")
			role, _ := cmd.Flags().GetString("role")
			assignments := []*stream.RoleAssignment{{ChannelRole: role, UserID: userID}}

			_, err = c.Channel(chType, chID).AssignRole(cmd.Context(), assignments, nil)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully assigned role [%s] to [%s].\n", role, userID)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("user-id", "u", "", "[required] User id to assign a role to")
	fl.StringP("role", "r", "", "[required] Channel role name to assign")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")
	_ = cmd.MarkFlagRequired("role")

	return cmd
}

func hideCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hide-channel --type [channel-type] --id [channel-id] --user-id [user-id]",
		Short: "Hide a channel",
		Long: heredoc.Doc(`
			Hiding a channel will remove it from query channel requests for that
			user until a new message is added. Please keep in mind that hiding a channel
			is only available to members of that channel.
			You can still retrieve the list of hidden channels using the { "hidden" : true } query parameter.
		`),
		Example: heredoc.Doc(`
			# Hide a 'red-team' channel for user 'joe'
			$ stream-cli chat hide-channel --type messaging --id red-team --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")
			userID, _ := cmd.Flags().GetString("user-id")

			_, err = c.Channel(chType, chID).Hide(cmd.Context(), userID)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully hid channel for " + userID + "\n")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("user-id", "u", "", "[required] User id to hide the channel to")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-channel --type [channel-type] --id [channel-id] --user-id [user-id]",
		Short: "Show a channel",
		Long: heredoc.Doc(`
			Hiding a channel will remove it from query channel requests for that
			user until a new message is added.
			As opposed to this, showing a channel will add it to query channel requests for that user.
		`),
		Example: heredoc.Doc(`
			# Show a 'red-team' channel for user 'joe'
			$ stream-cli chat show-channel --type messaging --id red-team --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")
			userID, _ := cmd.Flags().GetString("user-id")

			_, err = c.Channel(chType, chID).Show(cmd.Context(), userID)
			if err != nil {
				return err
			}

			cmd.Println("Successfully shown channel for " + userID)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("user-id", "u", "", "[required] User id to show the channel to")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func truncateChannelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "truncate-channel --type [channel-type] --id [channel-id] [flags]",
		Short: "Truncate a channel",
		Long: heredoc.Doc(`
			Truncates a channel by removing all messages but keeping the channel metadata and members.

			Optional flags allow you to perform a hard delete, add a system message, skip push notifications, 
			and define the truncating user ID (for server-side calls).
		`),
		Example: heredoc.Doc(`
			# Truncate messages in 'general' channel of type messaging
			$ stream-cli chat truncate-channel --type messaging --id general

			# Truncate with hard delete and system message
			$ stream-cli chat truncate-channel --type messaging --id general --hard --message "Channel reset" --user-id system-user
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			typ, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			hard, _ := cmd.Flags().GetBool("hard")
			userID, _ := cmd.Flags().GetString("user-id")
			messageUserID, _ := cmd.Flags().GetString("message-user-id")
			msgText, _ := cmd.Flags().GetString("message")
			skipPush, _ := cmd.Flags().GetBool("skip-push")

			// If a message is provided, a message user id must also be provided
			if msgText != "" && messageUserID == "" {
				return errors.New("when using --message, you must also supply --message-user-id")
			}

			var opts []stream.TruncateOption

			if hard {
				opts = append(opts, stream.TruncateWithHardDelete())
			}
			if skipPush {
				opts = append(opts, stream.TruncateWithSkipPush())
			}
			if userID != "" {
				opts = append(opts, stream.TruncateWithUserID(userID))
			}
			if msgText != "" {
				opts = append(opts, stream.TruncateWithMessage(&stream.Message{
					Text: msgText,
					User: &stream.User{ID: messageUserID},
				}))
			}

			ch := client.Channel(typ, id)
			_, err = ch.Truncate(cmd.Context(), opts...)
			if err != nil {
				return err
			}
			cmd.Printf("Successfully truncated channel [%s]\n", id)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging'")
	fl.StringP("id", "i", "", "[required] Channel ID")
	fl.String("user-id", "", "[optional] User ID who performs the truncation")
	fl.String("message", "", "[optional] System message to include in truncation (requires --message-user-id)")
	fl.String("message-user-id", "", "[optional] User id for the message to include in truncation (required if --message is set)")
	fl.Bool("hard", false, "[optional] Permanently delete messages instead of hiding them")
	fl.Bool("skip-push", false, "[optional] Skip push notifications")

	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
