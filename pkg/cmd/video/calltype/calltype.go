package calltype

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listCmd(),
		getCmd(),
		createCmd(),
		updateCmd(),
		deleteCmd(),
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

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-call-types",
		Short: "List all call types",
		Example: heredoc.Doc(`
			$ stream-cli video list-call-types
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "video/calltypes", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-call-type [name]",
		Short: "Get a call type by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "video/calltypes/"+args[0], nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-call-type --properties [json]",
		Short: "Create a call type",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "video/calltypes", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Call type properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-call-type --name [name] --properties [json]",
		Short: "Update a call type",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "video/calltypes/"+name, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Call type name")
	fl.StringP("properties", "p", "", "[required] Call type properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-call-type [name]",
		Short: "Delete a call type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "DELETE", "video/calltypes/"+args[0], nil)
			if err != nil {
				return err
			}
			cmd.Printf("Successfully deleted call type [%s]\n", args[0])
			return nil
		},
	}
	return cmd
}
