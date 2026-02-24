package feedback

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		queryUserFeedbackCmd(),
	}
}

func queryUserFeedbackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-user-feedback --properties [json]",
		Short: "Query user reported feedback for calls",
		RunE: func(cmd *cobra.Command, args []string) error {
			appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
			if err != nil {
				return err
			}
			h := utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL)

			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "video/call/feedback", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[optional] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}
