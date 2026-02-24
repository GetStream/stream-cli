package unread

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		unreadCountsCmd(),
		unreadCountsBatchCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func unreadCountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unread-counts --user-id [user-id]",
		Short: "Get unread counts for a user",
		Example: heredoc.Doc(`
			# Get unread counts for user 'joe'
			$ stream-cli chat unread-counts --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user-id")

			resp, err := h.DoRequest(cmd.Context(), "GET", "unread?user_id="+userID, nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("user-id", "u", "", "[required] User ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func unreadCountsBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unread-counts-batch [user-id-1] [user-id-2] ...",
		Short: "Get unread counts for multiple users",
		Example: heredoc.Doc(`
			# Get unread counts for multiple users
			$ stream-cli chat unread-counts-batch user-1 user-2 user-3
		`),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"user_ids": args,
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "unread_batch", body)
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
