package bookmark

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		queryBookmarksCmd(),
		queryBookmarkFoldersCmd(),
		updateBookmarkFolderCmd(),
		deleteBookmarkFolderCmd(),
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

func doAction(cmd *cobra.Command, method, path string, body interface{}, msg string) error {
	h, err := getHTTPClient(cmd)
	if err != nil {
		return err
	}
	_, err = h.DoRequest(cmd.Context(), method, path, body)
	if err != nil {
		return err
	}
	cmd.Println(msg)
	return nil
}

func queryBookmarksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-bookmarks --properties [json]",
		Short: "Query bookmarks with filters",
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
			return doJSON(cmd, "POST", "api/v2/feeds/bookmarks/query", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[optional] Query filters as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func queryBookmarkFoldersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-bookmark-folders --properties [json]",
		Short: "Query bookmark folders with filters",
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
			return doJSON(cmd, "POST", "api/v2/feeds/bookmark_folders/query", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[optional] Query filters as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func updateBookmarkFolderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-bookmark-folder --folder-id [id] --properties [json]",
		Short: "Update a bookmark folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			folderID, _ := cmd.Flags().GetString("folder-id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", "api/v2/feeds/bookmark_folders/"+folderID, body)
		},
	}
	fl := cmd.Flags()
	fl.String("folder-id", "", "[required] Bookmark folder ID")
	fl.StringP("properties", "p", "", "[required] Folder update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("folder-id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteBookmarkFolderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-bookmark-folder --folder-id [id]",
		Short: "Delete a bookmark folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			folderID, _ := cmd.Flags().GetString("folder-id")
			return doAction(cmd, "DELETE", "api/v2/feeds/bookmark_folders/"+folderID, nil, "Successfully deleted bookmark folder")
		},
	}
	fl := cmd.Flags()
	fl.String("folder-id", "", "[required] Bookmark folder ID")
	_ = cmd.MarkFlagRequired("folder-id")
	return cmd
}
