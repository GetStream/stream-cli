package storage

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listCmd(),
		createCmd(),
		updateCmd(),
		deleteCmd(),
		checkCmd(),
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
		Use:   "list-external-storage",
		Short: "List external storage configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "video/external_storage", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-external-storage --properties [json]",
		Short: "Create external storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "POST", "video/external_storage", body)
			if err != nil {
				return err
			}
			cmd.Println("Successfully created external storage")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Storage properties as JSON")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-external-storage --name [name] --properties [json]",
		Short: "Update external storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PUT", "video/external_storage/"+name, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Storage name")
	fl.StringP("properties", "p", "", "[required] Storage properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-external-storage [name]",
		Short: "Delete external storage",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "DELETE", "video/external_storage/"+args[0], nil)
			if err != nil {
				return err
			}
			cmd.Printf("Successfully deleted external storage [%s]\n", args[0])
			return nil
		},
	}
	return cmd
}

func checkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-external-storage [name]",
		Short: "Check external storage configuration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return doJSON(cmd, "GET", "video/external_storage/"+args[0]+"/check", nil)
		},
	}
	cmd.Flags().StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}
