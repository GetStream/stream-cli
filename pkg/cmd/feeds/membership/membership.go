package membership

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		createMembershipLevelCmd(),
		queryMembershipLevelsCmd(),
		updateMembershipLevelCmd(),
		deleteMembershipLevelCmd(),
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

func createMembershipLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-membership-level --properties [json]",
		Short: "Create a new membership level",
		Example: heredoc.Doc(`
			$ stream-cli feeds create-membership-level --properties '{"id":"premium","name":"Premium","priority":10,"tags":["premium-content"]}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/membership_levels", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Membership level properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func queryMembershipLevelsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-membership-levels --properties [json]",
		Short: "Query membership levels",
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
			return doJSON(cmd, "POST", "api/v2/feeds/membership_levels/query", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[optional] Query filters as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func updateMembershipLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-membership-level --id [id] --properties [json]",
		Short: "Update a membership level",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", "api/v2/feeds/membership_levels/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Membership level ID")
	fl.StringP("properties", "p", "", "[required] Membership level update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteMembershipLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-membership-level --id [id]",
		Short: "Delete a membership level",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			return doAction(cmd, "DELETE", "api/v2/feeds/membership_levels/"+id, nil, "Successfully deleted membership level")
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Membership level ID")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
