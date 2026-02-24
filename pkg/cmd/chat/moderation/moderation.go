package moderation

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		queryBannedUsersCmd(),
		queryMessageFlagsCmd(),
		getRateLimitsCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func queryBannedUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-banned-users --filter [raw-json]",
		Short: "Find and filter channel scoped or global user bans",
		Example: heredoc.Doc(`
			# Query all banned users
			$ stream-cli chat query-banned-users --filter '{"channel_cid":"messaging:general"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			filterStr, _ := cmd.Flags().GetString("filter")
			var filter interface{}
			if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
				return err
			}

			payload := map[string]interface{}{"filter_conditions": filter}
			payloadJSON, _ := json.Marshal(payload)

			resp, err := h.DoRequest(cmd.Context(), "GET", "query_banned_users?payload="+string(payloadJSON), nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("filter", "f", "{}", "[required] Filter conditions as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func queryMessageFlagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-message-flags --filter [raw-json]",
		Short: "Find and filter message flags",
		Example: heredoc.Doc(`
			# Query all message flags
			$ stream-cli chat query-message-flags --filter '{}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			filterStr, _ := cmd.Flags().GetString("filter")
			var filter interface{}
			if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
				return err
			}

			payload := map[string]interface{}{"filter_conditions": filter}
			payloadJSON, _ := json.Marshal(payload)

			resp, err := h.DoRequest(cmd.Context(), "GET", "moderation/flags/message?payload="+string(payloadJSON), nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("filter", "f", "{}", "[optional] Filter conditions as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func getRateLimitsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-rate-limits",
		Short: "Get rate limits usage and quotas",
		Example: heredoc.Doc(`
			# Get rate limits for server-side
			$ stream-cli chat get-rate-limits --server-side

			# Get rate limits for web
			$ stream-cli chat get-rate-limits --web
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			serverSide, _ := cmd.Flags().GetBool("server-side")
			web, _ := cmd.Flags().GetBool("web")
			android, _ := cmd.Flags().GetBool("android")
			ios, _ := cmd.Flags().GetBool("ios")

			path := "rate_limits?"
			if serverSide {
				path += "server_side=true&"
			}
			if web {
				path += "web=true&"
			}
			if android {
				path += "android=true&"
			}
			if ios {
				path += "ios=true&"
			}

			resp, err := h.DoRequest(cmd.Context(), "GET", path, nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.Bool("server-side", false, "[optional] Include server-side limits")
	fl.Bool("web", false, "[optional] Include web limits")
	fl.Bool("android", false, "[optional] Include Android limits")
	fl.Bool("ios", false, "[optional] Include iOS limits")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}
