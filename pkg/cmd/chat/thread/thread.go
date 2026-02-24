package thread

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		queryThreadsCmd(),
		getThreadCmd(),
		updateThreadPartialCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func queryThreadsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-threads --user-id [user-id]",
		Short: "Query threads for a user",
		Example: heredoc.Doc(`
			# Query threads for a user
			$ stream-cli chat query-threads --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user-id")
			limit, _ := cmd.Flags().GetInt("limit")

			body := map[string]interface{}{
				"user_id": userID,
				"limit":   limit,
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "threads", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("user-id", "u", "", "[required] User ID")
	fl.IntP("limit", "l", 10, "[optional] Number of threads to return")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func getThreadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-thread [message-id]",
		Short: "Get a thread by parent message ID",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			# Get thread details
			$ stream-cli chat get-thread parent-msg-id
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			resp, err := h.DoRequest(cmd.Context(), "GET", "threads/"+args[0], nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func updateThreadPartialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-thread-partial --message-id [id] --set [raw-json] --unset [fields]",
		Short: "Partially update a thread",
		Example: heredoc.Doc(`
			# Update thread title
			$ stream-cli chat update-thread-partial --message-id msg-1 --set '{"title":"New title"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			msgID, _ := cmd.Flags().GetString("message-id")
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

			resp, err := h.DoRequest(cmd.Context(), "PATCH", "threads/"+msgID, body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("message-id", "m", "", "[required] Parent message ID")
	fl.StringP("set", "s", "", "[optional] JSON of key-value pairs to set")
	fl.String("unset", "", "[optional] Comma separated list of fields to unset")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("message-id")

	return cmd
}
