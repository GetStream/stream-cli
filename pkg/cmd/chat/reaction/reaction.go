package reaction

import (
	stream_chat "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getCmd(),
		sendCmd(),
		deleteCmd()}
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-reactions [message-id]",
		Short: "Get reactions for a message",
		Example: heredoc.Doc(`
			# Get reactions for a [08f64828-3bba-42bd-8430-c26a3634ee5c] message
			$ stream-cli chat get-reactions 08f64828-3bba-42bd-8430-c26a3634ee5c --output-format json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			msgID := args[0]

			resp, err := c.Channel("", "").GetReactions(cmd.Context(), msgID, nil)
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp.Reactions)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}

func sendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-reaction --message-id [message-id] --user-id [user-id] --reaction-type [reaction-type]",
		Short: "Send a reaction to a message",
		Long: heredoc.Doc(`
			Stream Chat has built-in support for user Reactions. Common examples are
			likes, comments, loves, etc. Reactions can be customized so that you
			are able to use any type of reaction your application requires.
		`),
		Example: heredoc.Doc(`
			# Send a reaction to a [08f64828-3bba-42bd-8430-c26a3634ee5c] message
			$ stream-cli chat send-reaction --message-id 08f64828-3bba-42bd-8430-c26a3634ee5c --user-id 12345 --reaction-type like
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			reactionType, _ := cmd.Flags().GetString("reaction-type")
			msgID, _ := cmd.Flags().GetString("message-id")
			userID, _ := cmd.Flags().GetString("user-id")

			r := &stream_chat.Reaction{Type: reactionType}

			_, err = c.Channel("", "").SendReaction(cmd.Context(), r, msgID, userID)
			if err != nil {
				return err
			}

			cmd.Println("Successfully sent reaction")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("reaction-type", "r", "", "[required] The reaction type to send")
	fl.StringP("message-id", "m", "", "[required] The message id to send the reaction to")
	fl.StringP("user-id", "u", "", "[required] The user id of the user sending the reaction")
	_ = cmd.MarkFlagRequired("reaction-type")
	_ = cmd.MarkFlagRequired("message-id")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-reaction --message-id [message-id] --reaction-type [reaction-type] --user-id [user-id]",
		Short: "Delete a reaction from a message",
		Example: heredoc.Doc(`
			# Delete a reaction from [08f64828-3bba-42bd-8430-c26a3634ee5c] message
			$ stream-cli chat delete-reaction --message-id 08f64828-3bba-42bd-8430-c26a3634ee5c --reaction-type like --user-id 12345
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			reactionType, _ := cmd.Flags().GetString("reaction-type")
			userID, _ := cmd.Flags().GetString("user-id")
			msgID, _ := cmd.Flags().GetString("message-id")

			_, err = c.Channel("", "").DeleteReaction(cmd.Context(), msgID, reactionType, userID)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted reaction")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("reaction-type", "r", "", "[required] The reaction type to delete")
	fl.StringP("message-id", "m", "", "[required] The message id to delete the reaction from")
	fl.StringP("user-id", "u", "", "[required] The user id of the user deleting the reaction")
	cmd.MarkFlagRequired("reaction-type")
	cmd.MarkFlagRequired("message-id")
	cmd.MarkFlagRequired("user-id")

	return cmd
}
