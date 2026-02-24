package edge

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getEdgesCmd(),
	}
}

func getEdgesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-edges",
		Short: "Get available edges",
		RunE: func(cmd *cobra.Command, args []string) error {
			appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
			if err != nil {
				return err
			}
			h := utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL)

			resp, err := h.DoRequest(cmd.Context(), "GET", "video/edges", nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}
