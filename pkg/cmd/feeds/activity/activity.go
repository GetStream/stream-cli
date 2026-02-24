package activity

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
		addActivityCmd(),
		getActivityCmd(),
		updateActivityCmd(),
		updateActivityPartialCmd(),
		deleteActivityCmd(),
		restoreActivityCmd(),
		upsertActivitiesCmd(),
		deleteActivitiesCmd(),
		queryActivitiesCmd(),
		updateActivitiesPartialBatchCmd(),
		addActivityReactionCmd(),
		deleteActivityReactionCmd(),
		queryActivityReactionsCmd(),
		addActivityBookmarkCmd(),
		updateActivityBookmarkCmd(),
		deleteActivityBookmarkCmd(),
		activityFeedbackCmd(),
		castPollVoteCmd(),
		deletePollVoteCmd(),
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

func addActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-activity --properties [json]",
		Short: "Add a single activity",
		Long:  "Create a new activity in one or more feeds",
		Example: heredoc.Doc(`
			$ stream-cli feeds add-activity --properties '{"type":"post","feeds":["user:1"],"text":"Hello world"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/activities", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Activity properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func getActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-activity --id [id]",
		Short: "Get an activity by ID",
		Example: heredoc.Doc(`
			$ stream-cli feeds get-activity --id abc-123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "GET", "api/v2/feeds/activities/"+id, nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Activity ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func updateActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-activity --id [id] --properties [json]",
		Short: "Full update an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "api/v2/feeds/activities/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Activity ID")
	fl.StringP("properties", "p", "", "[required] Activity properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func updateActivityPartialCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-activity-partial --id [id] --properties [json]",
		Short: "Partial update an activity",
		Long:  "Update specific fields of an activity using set/unset operations",
		Example: heredoc.Doc(`
			$ stream-cli feeds update-activity-partial --id abc-123 --properties '{"set":{"text":"updated"}}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", "api/v2/feeds/activities/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Activity ID")
	fl.StringP("properties", "p", "", "[required] Partial update as JSON with set/unset")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-activity --id [id]",
		Short: "Delete an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			hard, _ := cmd.Flags().GetBool("hard")
			path := "api/v2/feeds/activities/" + id
			if hard {
				path += "?hard_delete=true"
			}
			return doAction(cmd, "DELETE", path, nil, "Successfully deleted activity")
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Activity ID")
	fl.Bool("hard", false, "[optional] Hard delete")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func restoreActivityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore-activity --id [id]",
		Short: "Restore a soft-deleted activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			return doJSON(cmd, "POST", "api/v2/feeds/activities/"+id+"/restore", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Activity ID")
	fl.StringP("properties", "p", "", "[optional] Restore properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func upsertActivitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-activities --properties [json]",
		Short: "Create or update multiple activities",
		Example: heredoc.Doc(`
			$ stream-cli feeds upsert-activities --properties '{"activities":[{"type":"post","feeds":["user:1"],"text":"Hello"}]}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/activities/batch", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Batch activities as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteActivitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-activities --properties [json]",
		Short: "Delete multiple activities",
		Example: heredoc.Doc(`
			$ stream-cli feeds delete-activities --properties '{"ids":["id1","id2"]}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/activities/delete", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Delete request as JSON with ids array")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func queryActivitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-activities --properties [json]",
		Short: "Query activities with filters",
		Example: heredoc.Doc(`
			$ stream-cli feeds query-activities --properties '{"filter":{"user_id":"user-1"},"limit":10}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/activities/query", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func updateActivitiesPartialBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-activities-partial-batch --properties [json]",
		Short: "Batch partial update of multiple activities",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", "api/v2/feeds/activities/batch/partial", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Batch partial update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func addActivityReactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-activity-reaction --activity-id [id] --properties [json]",
		Short: "Add a reaction to an activity",
		Example: heredoc.Doc(`
			$ stream-cli feeds add-activity-reaction --activity-id abc-123 --properties '{"type":"like"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			activityID, _ := cmd.Flags().GetString("activity-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/activities/"+activityID+"/reactions", body)
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	fl.StringP("properties", "p", "", "[required] Reaction properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("activity-id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteActivityReactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-activity-reaction --activity-id [id] --type [type]",
		Short: "Remove a reaction from an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			activityID, _ := cmd.Flags().GetString("activity-id")
			reactionType, _ := cmd.Flags().GetString("type")
			return doJSON(cmd, "DELETE", fmt.Sprintf("api/v2/feeds/activities/%s/reactions/%s", activityID, reactionType), nil)
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	fl.StringP("type", "t", "", "[required] Reaction type")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("activity-id")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func queryActivityReactionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-activity-reactions --activity-id [id] --properties [json]",
		Short: "Query reactions on an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			return doJSON(cmd, "POST", "api/v2/feeds/activities/"+activityID+"/reactions/query", body)
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	fl.StringP("properties", "p", "", "[optional] Query filters as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("activity-id")
	return cmd
}

func addActivityBookmarkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-activity-bookmark --activity-id [id] --properties [json]",
		Short: "Bookmark an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			return doJSON(cmd, "POST", "api/v2/feeds/activities/"+activityID+"/bookmarks", body)
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	fl.StringP("properties", "p", "", "[optional] Bookmark properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("activity-id")
	return cmd
}

func updateActivityBookmarkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-activity-bookmark --activity-id [id] --properties [json]",
		Short: "Update a bookmark on an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			activityID, _ := cmd.Flags().GetString("activity-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", "api/v2/feeds/activities/"+activityID+"/bookmarks", body)
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	fl.StringP("properties", "p", "", "[required] Bookmark update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("activity-id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteActivityBookmarkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-activity-bookmark --activity-id [id]",
		Short: "Remove a bookmark from an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			activityID, _ := cmd.Flags().GetString("activity-id")
			return doAction(cmd, "DELETE", "api/v2/feeds/activities/"+activityID+"/bookmarks", nil, "Successfully deleted bookmark")
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	_ = cmd.MarkFlagRequired("activity-id")
	return cmd
}

func activityFeedbackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activity-feedback --activity-id [id] --properties [json]",
		Short: "Provide feedback on an activity",
		Long:  "Submit feedback including show_more, show_less, or hide options",
		RunE: func(cmd *cobra.Command, args []string) error {
			activityID, _ := cmd.Flags().GetString("activity-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/activities/"+activityID+"/feedback", body)
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	fl.StringP("properties", "p", "", "[required] Feedback as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("activity-id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func castPollVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cast-poll-vote --activity-id [id] --poll-id [id] --properties [json]",
		Short: "Cast a vote on a poll",
		RunE: func(cmd *cobra.Command, args []string) error {
			activityID, _ := cmd.Flags().GetString("activity-id")
			pollID, _ := cmd.Flags().GetString("poll-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			return doJSON(cmd, "POST", fmt.Sprintf("api/v2/feeds/activities/%s/polls/%s/vote", activityID, pollID), body)
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	fl.String("poll-id", "", "[required] Poll ID")
	fl.StringP("properties", "p", "", "[optional] Vote properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("activity-id")
	_ = cmd.MarkFlagRequired("poll-id")
	return cmd
}

func deletePollVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-poll-vote --activity-id [id] --poll-id [id] --vote-id [id]",
		Short: "Delete a vote from a poll",
		RunE: func(cmd *cobra.Command, args []string) error {
			activityID, _ := cmd.Flags().GetString("activity-id")
			pollID, _ := cmd.Flags().GetString("poll-id")
			voteID, _ := cmd.Flags().GetString("vote-id")
			return doJSON(cmd, "DELETE", fmt.Sprintf("api/v2/feeds/activities/%s/polls/%s/vote/%s", activityID, pollID, voteID), nil)
		},
	}
	fl := cmd.Flags()
	fl.String("activity-id", "", "[required] Activity ID")
	fl.String("poll-id", "", "[required] Poll ID")
	fl.String("vote-id", "", "[required] Vote ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("activity-id")
	_ = cmd.MarkFlagRequired("poll-id")
	_ = cmd.MarkFlagRequired("vote-id")
	return cmd
}
