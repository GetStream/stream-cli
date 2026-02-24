package guest

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		createGuestCmd(),
	}
}

func createGuestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-guest --properties [json]",
		Short: "Create a guest user for video",
		RunE: func(cmd *cobra.Command, args []string) error {
			appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
			if err != nil {
				return err
			}
			h := utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL)

			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "video/guest", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Guest user properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}
