package visibility

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listFeedVisibilitiesCmd(),
		getFeedVisibilityCmd(),
		updateFeedVisibilityCmd(),
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

func listFeedVisibilitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-feed-visibilities",
		Short: "List all feed visibility configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "api/v2/feeds/feed_visibilities", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func getFeedVisibilityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-feed-visibility --name [name]",
		Short: "Get a feed visibility configuration",
		Long:  "Get visibility config by name: public, visible, followers, members, or private",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			return doJSON(cmd, "GET", "api/v2/feeds/feed_visibilities/"+name, nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Visibility name (public, visible, followers, members, private)")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func updateFeedVisibilityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-feed-visibility --name [name] --properties [json]",
		Short: "Update a feed visibility configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "api/v2/feeds/feed_visibilities/"+name, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Visibility name (public, visible, followers, members, private)")
	fl.StringP("properties", "p", "", "[required] Visibility update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}
