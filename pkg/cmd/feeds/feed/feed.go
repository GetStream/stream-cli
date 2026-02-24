package feed

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getOrCreateFeedCmd(),
		updateFeedCmd(),
		deleteFeedCmd(),
		createFeedsBatchCmd(),
		deleteFeedsBatchCmd(),
		queryFeedsCmd(),
		pinActivityCmd(),
		unpinActivityCmd(),
		markActivitiesCmd(),
		updateFeedMembersCmd(),
		queryFeedMembersCmd(),
		acceptFeedMemberInviteCmd(),
		rejectFeedMemberInviteCmd(),
		ownBatchCmd(),
		getFeedsRateLimitsCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
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

func feedPath(groupID, feedID string) string {
	return fmt.Sprintf("api/v2/feeds/feed_groups/%s/feeds/%s", groupID, feedID)
}

func getOrCreateFeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-or-create-feed --feed-group-id [group] --feed-id [id] --properties [json]",
		Short: "Get or create a feed",
		Example: heredoc.Doc(`
			$ stream-cli feeds get-or-create-feed --feed-group-id user --feed-id user-1 --properties '{}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			return doJSON(cmd, "POST", feedPath(groupID, feedID), body)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.StringP("properties", "p", "", "[optional] Feed properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	return cmd
}

func updateFeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-feed --feed-group-id [group] --feed-id [id] --properties [json]",
		Short: "Update a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", feedPath(groupID, feedID), body)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.StringP("properties", "p", "", "[required] Feed properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteFeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-feed --feed-group-id [group] --feed-id [id]",
		Short: "Delete a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			hard, _ := cmd.Flags().GetBool("hard")
			path := feedPath(groupID, feedID)
			if hard {
				path += "?hard_delete=true"
			}
			return doJSON(cmd, "DELETE", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.Bool("hard", false, "[optional] Hard delete")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	return cmd
}

func createFeedsBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-feeds-batch --properties [json]",
		Short: "Create multiple feeds at once",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/feeds/batch", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Batch feeds as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteFeedsBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-feeds-batch --properties [json]",
		Short: "Delete multiple feeds",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/feeds/delete", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Delete request as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func queryFeedsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-feeds --properties [json]",
		Short: "Query feeds with filters",
		Example: heredoc.Doc(`
			$ stream-cli feeds query-feeds --properties '{"filter":{"group_id":"user"},"limit":10}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/feeds/query", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func pinActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pin-activity --feed-group-id [group] --feed-id [id] --activity-id [id]",
		Short: "Pin an activity to a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			activityID, _ := cmd.Flags().GetString("activity-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			path := fmt.Sprintf("%s/activities/%s/pin", feedPath(groupID, feedID), activityID)
			return doJSON(cmd, "POST", path, body)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.String("activity-id", "", "[required] Activity ID to pin")
	fl.StringP("properties", "p", "", "[optional] Pin properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	_ = cmd.MarkFlagRequired("activity-id")
	return cmd
}

func unpinActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpin-activity --feed-group-id [group] --feed-id [id] --activity-id [id]",
		Short: "Unpin an activity from a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			activityID, _ := cmd.Flags().GetString("activity-id")
			path := fmt.Sprintf("%s/activities/%s/pin", feedPath(groupID, feedID), activityID)
			return doJSON(cmd, "DELETE", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.String("activity-id", "", "[required] Activity ID to unpin")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	_ = cmd.MarkFlagRequired("activity-id")
	return cmd
}

func markActivitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mark-activities --feed-group-id [group] --feed-id [id] --properties [json]",
		Short: "Mark activities as read/seen/watched",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			path := fmt.Sprintf("%s/activities/mark/batch", feedPath(groupID, feedID))
			return doAction(cmd, "POST", path, body, "Successfully marked activities")
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.StringP("properties", "p", "", "[required] Mark request as JSON")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func updateFeedMembersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-feed-members --feed-group-id [group] --feed-id [id] --properties [json]",
		Short: "Add, remove, or set members for a feed",
		Example: heredoc.Doc(`
			$ stream-cli feeds update-feed-members --feed-group-id user --feed-id user-1 --properties '{"operation":"upsert","members":[{"user_id":"user-2","role":"member"}]}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			path := fmt.Sprintf("%s/members", feedPath(groupID, feedID))
			return doJSON(cmd, "PATCH", path, body)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.StringP("properties", "p", "", "[required] Members update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func queryFeedMembersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-feed-members --feed-group-id [group] --feed-id [id] --properties [json]",
		Short: "Query members of a feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			path := fmt.Sprintf("%s/members/query", feedPath(groupID, feedID))
			return doJSON(cmd, "POST", path, body)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.StringP("properties", "p", "", "[optional] Query filters as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	return cmd
}

func acceptFeedMemberInviteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-feed-member-invite --feed-group-id [group] --feed-id [id]",
		Short: "Accept a feed member invite",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			path := fmt.Sprintf("%s/members/accept", feedPath(groupID, feedID))
			return doJSON(cmd, "POST", path, body)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.StringP("properties", "p", "", "[optional] Accept properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	return cmd
}

func rejectFeedMemberInviteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject-feed-member-invite --feed-group-id [group] --feed-id [id]",
		Short: "Reject a feed member invite",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			feedID, _ := cmd.Flags().GetString("feed-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			path := fmt.Sprintf("%s/members/reject", feedPath(groupID, feedID))
			return doJSON(cmd, "POST", path, body)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.String("feed-id", "", "[required] Feed ID")
	fl.StringP("properties", "p", "", "[optional] Reject properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	_ = cmd.MarkFlagRequired("feed-id")
	return cmd
}

func ownBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "own-batch --properties [json]",
		Short: "Get own fields for multiple feeds",
		Long:  "Retrieves own_follows, own_capabilities, and/or own_membership for multiple feeds",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/feeds/own/batch", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Own batch request as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func getFeedsRateLimitsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-feeds-rate-limits",
		Short: "Get rate limits for feeds operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "api/v2/feeds/feeds/rate_limits", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}
