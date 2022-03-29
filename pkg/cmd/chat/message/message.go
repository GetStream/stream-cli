package message

import (
	"strings"

	stream "github.com/GetStream/stream-chat-go/v5"
	chatUtils "github.com/GetStream/stream-cli/pkg/cmd/chat/utils"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		sendCmd(),
		getCmd(),
		getMultipleCmd(),
		partialUpdateCmd(),
		deleteCmd()}
}

func sendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-message --channel-type [channel-type] --channel-id [channel-id] --text [text] --user [user-id]",
		Short: "Send a message to a channel",
		Example: heredoc.Doc(`
			# Sends a message to 'redteam' channel of 'messaging' channel type
			$ stream-cli chat send-message --channel-type messaging --channel-id redteam --text "Hello World!" --user "user-1"

			# Sends a message to 'redteam' channel of 'livestream' channel type with an URL attachment
			$ stream-cli chat send-message --channel-type livestream --channel-id redteam --attachment "https://example.com/image.png" --text "Hello World!" --user "user-1"

			# You can also send a message with a local file attachment
			# In this scenario, we'll upload the file first then send the message
			$ stream-cli chat send-message --channel-type livestream --channel-id redteam --attachment "./image.png" --text "Hello World!" --user "user-1"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("channel-type")
			chId, _ := cmd.Flags().GetString("channel-id")
			user, _ := cmd.Flags().GetString("user")
			text, _ := cmd.Flags().GetString("text")
			attachment, _ := cmd.Flags().GetString("attachment")

			m := &stream.Message{Text: text}

			if attachment != "" {
				if strings.HasPrefix(attachment, "http") {
					m.Attachments = []*stream.Attachment{{AssetURL: attachment}}
				} else {
					uri, err := chatUtils.UploadFile(c, cmd, chType, chId, user, attachment, "")
					if err != nil {
						return err
					}
					m.Attachments = []*stream.Attachment{{AssetURL: uri}}
				}
			}

			msg, err := c.Channel(chType, chId).SendMessage(cmd.Context(), m, user)
			if err != nil {
				return err
			}

			cmd.Printf("Message successfully sent. Message id: [%s]\n", msg.Message.ID)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channel-type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("channel-id", "i", "", "[required] Channel id")
	fl.StringP("user", "u", "", "[required] User id")
	fl.String("text", "", "[required] Text of the message")
	fl.StringP("attachment", "a", "", "[optional] URL of the an attachment")
	cmd.MarkFlagRequired("channel-type")
	cmd.MarkFlagRequired("channel-id")
	cmd.MarkFlagRequired("user")
	cmd.MarkFlagRequired("text")

	return cmd
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-message [message-id] --output-format [json|tree]",
		Short: "Return a single message",
		Example: heredoc.Doc(`
			# Returns a message with id 'msgid-1'
			$ stream-cli chat get-message msgid-1
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			msg, err := c.GetMessage(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			utils.PrintObject(cmd, msg)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}

func getMultipleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-messages --channel-type [channel-type] --channel-id [channel-id] --output-format [json|tree] [message-id-1] [message-id-2] [message-id ...]",
		Short: "Return multiple messages",
		Example: heredoc.Doc(`
			# Returns 3 messages of 'redteam' channel of 'messaging' channel type
			$ stream-cli chat get-messages --channel-type messaging --channel-id redteam msgid-1 msgid-2 msgid-3
		`),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("channel-type")
			chId, _ := cmd.Flags().GetString("channel-id")

			messages, err := c.Channel(chType, chId).GetMessages(cmd.Context(), args)
			if err != nil {
				return err
			}

			utils.PrintObject(cmd, messages)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channel-type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("channel-id", "i", "", "[required] Channel id")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	cmd.MarkFlagRequired("channel-type")
	cmd.MarkFlagRequired("channel-id")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-message [message-id]",
		Short: "Delete a message",
		Long: heredoc.Doc(`
			You can delete a message by calling DeleteMessage and including a message
			with an ID. Messages can be soft deleted or hard deleted. Unless specified
			via the hard parameter, messages are soft deleted. Be aware that deleting
			a message doesn't delete its attachments. 
		`),
		Example: heredoc.Doc(`
			# Soft deletes a message with id 'msgid-1'
			$ stream-cli chat delete-message msgid-1

			# Hard deletes a message with id 'msgid-2'
			$ stream-cli chat delete-message msgid-2 --hard
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			hard, _ := cmd.Flags().GetBool("hard")

			if hard {
				_, err = c.HardDeleteMessage(cmd.Context(), args[0])
			} else {
				_, err = c.DeleteMessage(cmd.Context(), args[0])

			}

			if err != nil {
				return err
			}

			cmd.Printf("Message successfully deleted.\n")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.BoolP("hard", "H", false, "[optional] Hard delete message. Default is false")

	return cmd
}

func partialUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-message-partial --message-id [message-id] --user [user-id] --set [key-value-pairs] --unset [property-names]",
		Short: "Partially update a message",
		Long: heredoc.Doc(`
			A partial update can be used to set and unset specific fields when it
			is necessary to retain additional data fields on the object. AKA a patch style update.
		`),
		Example: heredoc.Doc(`
			# Partially updates a message with id 'msgid-1'. Updates a custom field and removes the silent flag.
			$ stream-cli chat update-message-partial -message-id msgid-1 --set importance=low --unset silent
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			msgId, _ := cmd.Flags().GetString("message-id")
			user, _ := cmd.Flags().GetString("user")
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

			_, err = c.PartialUpdateMessage(cmd.Context(), msgId, &stream.MessagePartialUpdateRequest{
				UserID: user,
				PartialUpdate: stream.PartialUpdate{
					Set:   s,
					Unset: u,
				},
			})
			if err != nil {
				return err
			}

			cmd.Println("Successfully updated message.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("message-id", "m", "", "[required] Message id")
	fl.StringP("user", "u", "", "[required] User id")
	fl.StringToStringP("set", "s", map[string]string{}, "[optional] Comma-separated key-value pairs to set")
	fl.String("unset", "", "[optional] Comma separated list of properties to unset")
	cmd.MarkFlagRequired("message-id")

	return cmd
}
