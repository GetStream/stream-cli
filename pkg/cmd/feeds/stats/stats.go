package stats

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		queryUsageStatsCmd(),
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

func queryUsageStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-usage-stats --properties [json]",
		Short: "Query feeds usage statistics",
		Long:  "Retrieve usage statistics including activity count, follow count, and API request count",
		Example: heredoc.Doc(`
			$ stream-cli feeds query-usage-stats --properties '{"from":"2025-01-01","to":"2025-01-31"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			return doJSON(cmd, "POST", "api/v2/feeds/stats/usage", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[optional] Stats query as JSON (from, to dates)")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}
