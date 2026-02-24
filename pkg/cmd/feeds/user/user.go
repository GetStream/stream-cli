package user

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		deleteFeedUserDataCmd(),
		exportFeedUserDataCmd(),
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

func deleteFeedUserDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-feed-user-data --user-id [id]",
		Short: "Delete all feed data for a user",
		Long:  "Delete feeds, activities, follows, comments, reactions, bookmarks, and collections owned by the user",
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, _ := cmd.Flags().GetString("user-id")
			hard, _ := cmd.Flags().GetBool("hard")
			body := map[string]interface{}{}
			if hard {
				body["hard_delete"] = true
			}
			return doJSON(cmd, "POST", "api/v2/feeds/users/"+userID+"/delete", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("user-id", "u", "", "[required] User ID")
	fl.Bool("hard", false, "[optional] Hard delete instead of soft delete")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("user-id")
	return cmd
}

func exportFeedUserDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-feed-user-data --user-id [id]",
		Short: "Export all feed data for a user",
		Long:  "Export user profile, feeds, activities, follows, comments, reactions, bookmarks, and collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			userID, _ := cmd.Flags().GetString("user-id")
			return doJSON(cmd, "POST", "api/v2/feeds/users/"+userID+"/export", map[string]interface{}{})
		},
	}
	fl := cmd.Flags()
	fl.StringP("user-id", "u", "", "[required] User ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("user-id")
	return cmd
}
