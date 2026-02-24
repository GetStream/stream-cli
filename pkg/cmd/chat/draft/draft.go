package draft

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getDraftCmd(),
		deleteDraftCmd(),
		queryDraftsCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func getDraftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-draft --type [channel-type] --id [channel-id] --user-id [user-id]",
		Short: "Get a draft message",
		Example: heredoc.Doc(`
			# Get a draft for a channel
			$ stream-cli chat get-draft --type messaging --id general --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")
			userID, _ := cmd.Flags().GetString("user-id")

			path := "channels/" + chType + "/" + chID + "/draft?user_id=" + userID

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
	fl.StringP("type", "t", "", "[required] Channel type")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("user-id", "u", "", "[required] User id")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func deleteDraftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-draft --type [channel-type] --id [channel-id] --user-id [user-id]",
		Short: "Delete a draft message",
		Example: heredoc.Doc(`
			# Delete a draft
			$ stream-cli chat delete-draft --type messaging --id general --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			chType, _ := cmd.Flags().GetString("type")
			chID, _ := cmd.Flags().GetString("id")
			userID, _ := cmd.Flags().GetString("user-id")

			path := "channels/" + chType + "/" + chID + "/draft?user_id=" + userID

			_, err = h.DoRequest(cmd.Context(), "DELETE", path, nil)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted draft")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("user-id", "u", "", "[required] User id")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func queryDraftsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-drafts --user-id [user-id]",
		Short: "Query draft messages for a user",
		Example: heredoc.Doc(`
			# Query drafts for a user
			$ stream-cli chat query-drafts --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user-id")
			body := map[string]interface{}{
				"user_id": userID,
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "drafts/query", body)
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
