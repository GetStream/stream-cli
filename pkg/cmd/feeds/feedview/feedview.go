package feedview

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listFeedViewsCmd(),
		createFeedViewCmd(),
		getFeedViewCmd(),
		getOrCreateFeedViewCmd(),
		updateFeedViewCmd(),
		deleteFeedViewCmd(),
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

func listFeedViewsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-feed-views",
		Short: "List all feed views",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "api/v2/feeds/feed_views", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func createFeedViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-feed-view --properties [json]",
		Short: "Create a custom feed view",
		Example: heredoc.Doc(`
			$ stream-cli feeds create-feed-view --properties '{"id":"trending","activity_selectors":[{"type":"popular"}]}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/feed_views", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Feed view properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func getFeedViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-feed-view --id [id]",
		Short: "Get a feed view by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "GET", "api/v2/feeds/feed_views/"+id, nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Feed view ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func getOrCreateFeedViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-or-create-feed-view --id [id] --properties [json]",
		Short: "Get or create a feed view",
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
			return doJSON(cmd, "POST", "api/v2/feeds/feed_views/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Feed view ID")
	fl.StringP("properties", "p", "", "[optional] Feed view properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func updateFeedViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-feed-view --id [id] --properties [json]",
		Short: "Update a feed view",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "api/v2/feeds/feed_views/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Feed view ID")
	fl.StringP("properties", "p", "", "[required] Feed view properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteFeedViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-feed-view --id [id]",
		Short: "Delete a feed view",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			return doAction(cmd, "DELETE", "api/v2/feeds/feed_views/"+id, nil, "Successfully deleted feed view")
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Feed view ID")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
