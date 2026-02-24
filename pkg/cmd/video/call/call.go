package call

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getCallCmd(),
		updateCallCmd(),
		getOrCreateCallCmd(),
		deleteCallCmd(),
		queryCallsCmd(),
		getActiveCallsStatusCmd(),
		goLiveCmd(),
		stopLiveCmd(),
		endCallCmd(),
		ringCallCmd(),
		blockUserCmd(),
		unblockUserCmd(),
		kickUserCmd(),
		muteUsersCmd(),
		sendEventCmd(),
		pinCmd(),
		unpinCmd(),
		updateCallMembersCmd(),
		queryCallMembersCmd(),
		updateUserPermissionsCmd(),
		sendClosedCaptionCmd(),
		collectUserFeedbackCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func callPath(t, id string) string {
	return "video/call/" + t + "/" + id
}

func doJSON(cmd *cobra.Command, method, path string, body interface{}) error {
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

func doAction(cmd *cobra.Command, method, path string, body interface{}, msg string) error {
	h, err := getHTTPClient(cmd)
	if err != nil {
		return err
	}
	_, err = h.DoRequest(cmd.Context(), method, path, body)
	if err != nil {
		return err
	}
	cmd.Println(msg)
	return nil
}

func getCallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-call --type [type] --id [id]",
		Short: "Get call details",
		Example: heredoc.Doc(`
			$ stream-cli video get-call --type default --id call-123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "GET", callPath(t, id), nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func updateCallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-call --type [type] --id [id] --properties [json]",
		Short: "Update a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", callPath(t, id), body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Call properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func getOrCreateCallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-or-create-call --type [type] --id [id] --properties [json]",
		Short: "Get or create a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			}
			return doJSON(cmd, "POST", callPath(t, id), body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[optional] Call properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func deleteCallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-call --type [type] --id [id]",
		Short: "Delete a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			hard, _ := cmd.Flags().GetBool("hard")
			body := map[string]interface{}{}
			if hard {
				body["hard"] = true
			}
			return doAction(cmd, "POST", callPath(t, id)+"/delete", body, "Successfully deleted call")
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.Bool("hard", false, "[optional] Hard delete")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func queryCallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-calls --filter [json]",
		Short: "Query calls",
		RunE: func(cmd *cobra.Command, args []string) error {
			filterStr, _ := cmd.Flags().GetString("filter")
			body := map[string]interface{}{}
			if filterStr != "" {
				var f interface{}
				if err := json.Unmarshal([]byte(filterStr), &f); err != nil {
					return err
				}
				body["filter_conditions"] = f
			}
			limit, _ := cmd.Flags().GetInt("limit")
			if limit > 0 {
				body["limit"] = limit
			}
			return doJSON(cmd, "POST", "video/calls", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("filter", "f", "", "[optional] Filter conditions as JSON")
	fl.IntP("limit", "l", 0, "[optional] Limit results")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func getActiveCallsStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-active-calls-status",
		Short: "Get status of all active calls",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "video/active_calls_status", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func goLiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go-live --type [type] --id [id]",
		Short: "Set call as live",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "POST", callPath(t, id)+"/go_live", map[string]interface{}{})
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func stopLiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-live --type [type] --id [id]",
		Short: "Set call as not live",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "POST", callPath(t, id)+"/stop_live", map[string]interface{}{})
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func endCallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "end-call --type [type] --id [id]",
		Short: "End a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			return doAction(cmd, "POST", callPath(t, id)+"/mark_ended", map[string]interface{}{}, "Successfully ended call")
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func ringCallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ring-call --type [type] --id [id]",
		Short: "Ring call users",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "POST", callPath(t, id)+"/ring", map[string]interface{}{})
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func blockUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-user --type [type] --id [id] --user-id [user-id]",
		Short: "Block a user on a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			uid, _ := cmd.Flags().GetString("user-id")
			return doAction(cmd, "POST", callPath(t, id)+"/block", map[string]interface{}{"user_id": uid}, "Successfully blocked user")
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("user-id", "u", "", "[required] User ID to block")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")
	return cmd
}

func unblockUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unblock-user --type [type] --id [id] --user-id [user-id]",
		Short: "Unblock a user on a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			uid, _ := cmd.Flags().GetString("user-id")
			return doAction(cmd, "POST", callPath(t, id)+"/unblock", map[string]interface{}{"user_id": uid}, "Successfully unblocked user")
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("user-id", "u", "", "[required] User ID to unblock")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")
	return cmd
}

func kickUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kick-user --type [type] --id [id] --user-id [user-id]",
		Short: "Kick a user from a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			uid, _ := cmd.Flags().GetString("user-id")
			return doAction(cmd, "POST", callPath(t, id)+"/kick", map[string]interface{}{"user_id": uid}, "Successfully kicked user")
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("user-id", "u", "", "[required] User ID to kick")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")
	return cmd
}

func muteUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mute-users --type [type] --id [id] --properties [json]",
		Short: "Mute users on a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", callPath(t, id)+"/mute_users", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Mute properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func sendEventCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-event --type [type] --id [id] --properties [json]",
		Short: "Send a custom event to a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", callPath(t, id)+"/event", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Event properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func pinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin --type [type] --id [id] --properties [json]",
		Short: "Pin a participant on a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", callPath(t, id)+"/pin", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Pin properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func unpinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpin --type [type] --id [id] --properties [json]",
		Short: "Unpin a participant on a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", callPath(t, id)+"/unpin", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Unpin properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func updateCallMembersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-call-members --type [type] --id [id] --properties [json]",
		Short: "Update call members",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", callPath(t, id)+"/members", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Members update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func queryCallMembersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-call-members --properties [json]",
		Short: "Query call members",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "video/call/members", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func updateUserPermissionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-user-permissions --type [type] --id [id] --properties [json]",
		Short: "Update user permissions on a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", callPath(t, id)+"/user_permissions", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Permissions as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func sendClosedCaptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-closed-caption --type [type] --id [id] --properties [json]",
		Short: "Send a closed caption",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", callPath(t, id)+"/closed_captions", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Caption data as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func collectUserFeedbackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect-user-feedback --type [type] --id [id] --properties [json]",
		Short: "Collect user feedback for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", callPath(t, id)+"/feedback", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[required] Feedback data as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}
