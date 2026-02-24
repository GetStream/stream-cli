package og

import (
	"encoding/json"
	"net/url"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getOGCmd(),
	}
}

func getOGCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-og --url [url] --output-format [json|tree]",
		Short: "Get OpenGraph attachment for a URL",
		Example: heredoc.Doc(`
			# Get OG data for a URL
			$ stream-cli chat get-og --url "https://getstream.io"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
			if err != nil {
				return err
			}
			h := utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL)

			ogURL, _ := cmd.Flags().GetString("url")

			resp, err := h.DoRequest(cmd.Context(), "GET", "og?url="+url.QueryEscape(ogURL), nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.String("url", "", "[required] URL to scrape for OpenGraph data")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	_ = cmd.MarkFlagRequired("url")

	return cmd
}
