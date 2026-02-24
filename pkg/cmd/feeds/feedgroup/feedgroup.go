package feedgroup

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
		listFeedGroupsCmd(),
		createFeedGroupCmd(),
		getFeedGroupCmd(),
		getOrCreateFeedGroupCmd(),
		updateFeedGroupCmd(),
		deleteFeedGroupCmd(),
		getFollowSuggestionsCmd(),
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

func listFeedGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-feed-groups",
		Short: "List all feed groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "api/v2/feeds/feed_groups", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func createFeedGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-feed-group --properties [json]",
		Short: "Create a new feed group",
		Example: heredoc.Doc(`
			$ stream-cli feeds create-feed-group --properties '{"id":"timeline"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/feed_groups", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Feed group properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func getFeedGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-feed-group --id [id]",
		Short: "Get a feed group by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "GET", "api/v2/feeds/feed_groups/"+id, nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Feed group ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func getOrCreateFeedGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-or-create-feed-group --id [id] --properties [json]",
		Short: "Get or create a feed group",
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
			return doJSON(cmd, "POST", "api/v2/feeds/feed_groups/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Feed group ID")
	fl.StringP("properties", "p", "", "[optional] Feed group properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func updateFeedGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-feed-group --id [id] --properties [json]",
		Short: "Update a feed group",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "api/v2/feeds/feed_groups/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Feed group ID")
	fl.StringP("properties", "p", "", "[required] Feed group properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteFeedGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-feed-group --id [id]",
		Short: "Delete a feed group",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			hard, _ := cmd.Flags().GetBool("hard")
			path := "api/v2/feeds/feed_groups/" + id
			if hard {
				path += "?hard_delete=true"
			}
			return doAction(cmd, "DELETE", path, nil, "Successfully deleted feed group")
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Feed group ID")
	fl.Bool("hard", false, "[optional] Hard delete")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func getFollowSuggestionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-follow-suggestions --feed-group-id [id]",
		Short: "Get follow suggestions for a feed group",
		RunE: func(cmd *cobra.Command, args []string) error {
			groupID, _ := cmd.Flags().GetString("feed-group-id")
			path := fmt.Sprintf("api/v2/feeds/feed_groups/%s/follow_suggestions", groupID)
			limit, _ := cmd.Flags().GetInt("limit")
			if limit > 0 {
				path += fmt.Sprintf("?limit=%d", limit)
			}
			return doJSON(cmd, "GET", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.String("feed-group-id", "", "[required] Feed group ID")
	fl.IntP("limit", "l", 0, "[optional] Limit results")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("feed-group-id")
	return cmd
}
