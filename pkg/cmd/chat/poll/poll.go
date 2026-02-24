package poll

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		createPollCmd(),
		getPollCmd(),
		updatePollCmd(),
		updatePollPartialCmd(),
		deletePollCmd(),
		queryPollsCmd(),
		createPollOptionCmd(),
		updatePollOptionCmd(),
		getPollOptionCmd(),
		deletePollOptionCmd(),
		castPollVoteCmd(),
		deletePollVoteCmd(),
		queryPollVotesCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func doAndPrint(cmd *cobra.Command, method, path string, body interface{}) error {
	h, err := getHTTPClient(cmd)
	if err != nil {
		return err
	}

	resp, err := h.DoRequest(cmd.Context(), method, path, body)
	if err != nil {
		return err
	}

	var result interface{}
	_ = json.Unmarshal(resp, &result)
	return utils.PrintObject(cmd, result)
}

func createPollCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-poll --properties [raw-json]",
		Short: "Create a new poll",
		Example: heredoc.Doc(`
			# Create a poll
			$ stream-cli chat create-poll --properties '{"name":"Favorite color","options":[{"text":"Red"},{"text":"Blue"}]}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doAndPrint(cmd, "POST", "polls", body)
		},
	}

	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Poll properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}

func getPollCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-poll [poll-id]",
		Short: "Get a poll by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doAndPrint(cmd, "GET", "polls/"+args[0], nil)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func updatePollCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-poll --properties [raw-json]",
		Short: "Update a poll",
		Example: heredoc.Doc(`
			# Update a poll
			$ stream-cli chat update-poll --properties '{"id":"poll-id","name":"Updated name"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doAndPrint(cmd, "PUT", "polls", body)
		},
	}

	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Poll properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}

func updatePollPartialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-poll-partial --poll-id [id] --set [raw-json] --unset [fields]",
		Short: "Partially update a poll",
		Example: heredoc.Doc(`
			# Partially update a poll
			$ stream-cli chat update-poll-partial --poll-id my-poll --set '{"name":"New name"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			pollID, _ := cmd.Flags().GetString("poll-id")
			setStr, _ := cmd.Flags().GetString("set")
			unsetStr, _ := cmd.Flags().GetString("unset")

			body := map[string]interface{}{}
			if setStr != "" {
				var setMap map[string]interface{}
				if err := json.Unmarshal([]byte(setStr), &setMap); err != nil {
					return err
				}
				body["set"] = setMap
			}
			if unsetStr != "" {
				unset, _ := utils.GetStringSliceParam(cmd.Flags(), "unset")
				body["unset"] = unset
			}

			return doAndPrint(cmd, "PATCH", "polls/"+pollID, body)
		},
	}

	fl := cmd.Flags()
	fl.String("poll-id", "", "[required] Poll ID")
	fl.StringP("set", "s", "", "[optional] JSON of key-value pairs to set")
	fl.String("unset", "", "[optional] Comma separated list of fields to unset")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("poll-id")

	return cmd
}

func deletePollCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-poll [poll-id]",
		Short: "Delete a poll",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "DELETE", "polls/"+args[0], nil)
			if err != nil {
				return err
			}
			cmd.Printf("Successfully deleted poll [%s]\n", args[0])
			return nil
		},
	}

	return cmd
}

func queryPollsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-polls --filter [raw-json]",
		Short: "Query polls",
		RunE: func(cmd *cobra.Command, args []string) error {
			filterStr, _ := cmd.Flags().GetString("filter")
			body := map[string]interface{}{}
			if filterStr != "" {
				var filter interface{}
				if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
					return err
				}
				body["filter"] = filter
			}
			return doAndPrint(cmd, "POST", "polls/query", body)
		},
	}

	fl := cmd.Flags()
	fl.StringP("filter", "f", "", "[optional] Filter as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func createPollOptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-poll-option --poll-id [id] --text [text]",
		Short: "Create a poll option",
		RunE: func(cmd *cobra.Command, args []string) error {
			pollID, _ := cmd.Flags().GetString("poll-id")
			text, _ := cmd.Flags().GetString("text")
			body := map[string]interface{}{"text": text}
			return doAndPrint(cmd, "POST", "polls/"+pollID+"/options", body)
		},
	}

	fl := cmd.Flags()
	fl.String("poll-id", "", "[required] Poll ID")
	fl.String("text", "", "[required] Option text")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("poll-id")
	_ = cmd.MarkFlagRequired("text")

	return cmd
}

func updatePollOptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-poll-option --poll-id [id] --option-id [option-id] --text [text]",
		Short: "Update a poll option",
		RunE: func(cmd *cobra.Command, args []string) error {
			pollID, _ := cmd.Flags().GetString("poll-id")
			optionID, _ := cmd.Flags().GetString("option-id")
			text, _ := cmd.Flags().GetString("text")
			body := map[string]interface{}{"id": optionID, "text": text}
			return doAndPrint(cmd, "PUT", "polls/"+pollID+"/options", body)
		},
	}

	fl := cmd.Flags()
	fl.String("poll-id", "", "[required] Poll ID")
	fl.String("option-id", "", "[required] Option ID")
	fl.String("text", "", "[required] Option text")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("poll-id")
	_ = cmd.MarkFlagRequired("option-id")
	_ = cmd.MarkFlagRequired("text")

	return cmd
}

func getPollOptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-poll-option --poll-id [id] --option-id [option-id]",
		Short: "Get a poll option",
		RunE: func(cmd *cobra.Command, args []string) error {
			pollID, _ := cmd.Flags().GetString("poll-id")
			optionID, _ := cmd.Flags().GetString("option-id")
			return doAndPrint(cmd, "GET", "polls/"+pollID+"/options/"+optionID, nil)
		},
	}

	fl := cmd.Flags()
	fl.String("poll-id", "", "[required] Poll ID")
	fl.String("option-id", "", "[required] Option ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("poll-id")
	_ = cmd.MarkFlagRequired("option-id")

	return cmd
}

func deletePollOptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-poll-option --poll-id [id] --option-id [option-id]",
		Short: "Delete a poll option",
		RunE: func(cmd *cobra.Command, args []string) error {
			pollID, _ := cmd.Flags().GetString("poll-id")
			optionID, _ := cmd.Flags().GetString("option-id")

			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "DELETE", "polls/"+pollID+"/options/"+optionID, nil)
			if err != nil {
				return err
			}
			cmd.Println("Successfully deleted poll option")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.String("poll-id", "", "[required] Poll ID")
	fl.String("option-id", "", "[required] Option ID")
	_ = cmd.MarkFlagRequired("poll-id")
	_ = cmd.MarkFlagRequired("option-id")

	return cmd
}

func castPollVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cast-poll-vote --message-id [id] --poll-id [id] --option-id [id]",
		Short: "Cast a vote on a poll",
		RunE: func(cmd *cobra.Command, args []string) error {
			msgID, _ := cmd.Flags().GetString("message-id")
			pollID, _ := cmd.Flags().GetString("poll-id")
			optionID, _ := cmd.Flags().GetString("option-id")
			body := map[string]interface{}{
				"vote": map[string]interface{}{"option_id": optionID},
			}
			return doAndPrint(cmd, "POST", "messages/"+msgID+"/polls/"+pollID+"/vote", body)
		},
	}

	fl := cmd.Flags()
	fl.String("message-id", "", "[required] Message ID")
	fl.String("poll-id", "", "[required] Poll ID")
	fl.String("option-id", "", "[required] Option ID to vote for")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("message-id")
	_ = cmd.MarkFlagRequired("poll-id")
	_ = cmd.MarkFlagRequired("option-id")

	return cmd
}

func deletePollVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-poll-vote --message-id [id] --poll-id [id] --vote-id [id]",
		Short: "Delete a poll vote",
		RunE: func(cmd *cobra.Command, args []string) error {
			msgID, _ := cmd.Flags().GetString("message-id")
			pollID, _ := cmd.Flags().GetString("poll-id")
			voteID, _ := cmd.Flags().GetString("vote-id")

			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "DELETE", "messages/"+msgID+"/polls/"+pollID+"/vote/"+voteID, nil)
			if err != nil {
				return err
			}
			cmd.Println("Successfully deleted poll vote")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.String("message-id", "", "[required] Message ID")
	fl.String("poll-id", "", "[required] Poll ID")
	fl.String("vote-id", "", "[required] Vote ID")
	_ = cmd.MarkFlagRequired("message-id")
	_ = cmd.MarkFlagRequired("poll-id")
	_ = cmd.MarkFlagRequired("vote-id")

	return cmd
}

func queryPollVotesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-poll-votes --poll-id [id] --filter [raw-json]",
		Short: "Query poll votes",
		RunE: func(cmd *cobra.Command, args []string) error {
			pollID, _ := cmd.Flags().GetString("poll-id")
			filterStr, _ := cmd.Flags().GetString("filter")
			body := map[string]interface{}{}
			if filterStr != "" {
				var filter interface{}
				if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
					return err
				}
				body["filter"] = filter
			}
			return doAndPrint(cmd, "POST", "polls/"+pollID+"/votes", body)
		},
	}

	fl := cmd.Flags()
	fl.String("poll-id", "", "[required] Poll ID")
	fl.StringP("filter", "f", "", "[optional] Filter as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("poll-id")

	return cmd
}
