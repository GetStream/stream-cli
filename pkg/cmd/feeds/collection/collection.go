package collection

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		readCollectionsCmd(),
		createCollectionsCmd(),
		upsertCollectionsCmd(),
		updateCollectionsCmd(),
		deleteCollectionsCmd(),
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

func readCollectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read-collections --refs [refs]",
		Short: "Read collections by references",
		Example: heredoc.Doc(`
			$ stream-cli feeds read-collections --refs "food:pizza,food:burger"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			refs, _ := cmd.Flags().GetString("refs")
			path := "api/v2/feeds/collections?collection_refs=" + refs
			return doJSON(cmd, "GET", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.String("refs", "", "[required] Collection references (comma-separated, format: name:id)")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("refs")
	return cmd
}

func createCollectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-collections --properties [json]",
		Short: "Create multiple collections",
		Example: heredoc.Doc(`
			$ stream-cli feeds create-collections --properties '{"collections":[{"name":"food","custom":{"type":"pizza"}}]}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/collections", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Collections as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func upsertCollectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-collections --properties [json]",
		Short: "Create or update multiple collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "api/v2/feeds/collections", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Collections as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func updateCollectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-collections --properties [json]",
		Short: "Update multiple collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", "api/v2/feeds/collections", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Collection updates as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteCollectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-collections --refs [refs]",
		Short: "Delete multiple collections",
		Example: heredoc.Doc(`
			$ stream-cli feeds delete-collections --refs "food:pizza,food:burger"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			refs, _ := cmd.Flags().GetString("refs")
			path := "api/v2/feeds/collections?collection_refs=" + refs
			return doAction(cmd, "DELETE", path, nil, "Successfully deleted collections")
		},
	}
	fl := cmd.Flags()
	fl.String("refs", "", "[required] Collection references (comma-separated, format: name:id)")
	_ = cmd.MarkFlagRequired("refs")
	return cmd
}
