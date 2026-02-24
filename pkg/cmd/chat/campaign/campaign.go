package campaign

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getCampaignCmd(),
		queryCampaignsCmd(),
		startCampaignCmd(),
		stopCampaignCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func getCampaignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-campaign [campaign-id]",
		Short: "Get a campaign by ID",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			# Get campaign details
			$ stream-cli chat get-campaign campaign-123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			resp, err := h.DoRequest(cmd.Context(), "GET", "campaigns/"+args[0], nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func queryCampaignsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-campaigns --filter [raw-json]",
		Short: "Query campaigns",
		Example: heredoc.Doc(`
			# Query all campaigns
			$ stream-cli chat query-campaigns --filter '{}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			filterStr, _ := cmd.Flags().GetString("filter")
			body := map[string]interface{}{}
			if filterStr != "" {
				var filter interface{}
				if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
					return err
				}
				body["filter"] = filter
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "campaigns/query", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("filter", "f", "", "[optional] Filter conditions as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func startCampaignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-campaign [campaign-id]",
		Short: "Start or schedule a campaign",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			# Start a campaign immediately
			$ stream-cli chat start-campaign campaign-123

			# Schedule a campaign
			$ stream-cli chat start-campaign campaign-123 --scheduled-for "2025-12-01T10:00:00Z"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			body := map[string]interface{}{}
			scheduledFor, _ := cmd.Flags().GetString("scheduled-for")
			if scheduledFor != "" {
				body["scheduled_for"] = scheduledFor
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "campaigns/"+args[0]+"/start", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.String("scheduled-for", "", "[optional] Schedule time in RFC3339 format")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func stopCampaignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-campaign [campaign-id]",
		Short: "Stop a running campaign",
		Args:  cobra.ExactArgs(1),
		Example: heredoc.Doc(`
			# Stop a campaign
			$ stream-cli chat stop-campaign campaign-123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "campaigns/"+args[0]+"/stop", nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}
