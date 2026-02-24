package storage

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

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-external-storage --output-format [json|tree]",
		Short: "List all external storage configurations",
		Example: heredoc.Doc(`
			# List all external storage
			$ stream-cli chat list-external-storage
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			resp, err := h.DoRequest(cmd.Context(), "GET", "external_storage", nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-external-storage --properties [raw-json]",
		Short: "Create a new external storage",
		Example: heredoc.Doc(`
			# Create an S3 external storage
			$ stream-cli chat create-external-storage --properties '{"name":"my-s3","storage_type":"s3","bucket":"my-bucket","aws_s3":{"s3_region":"us-east-1"}}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}

			_, err = h.DoRequest(cmd.Context(), "POST", "external_storage", body)
			if err != nil {
				return err
			}

			cmd.Println("Successfully created external storage")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] JSON properties for the external storage")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-external-storage --name [name] --properties [raw-json]",
		Short: "Update an external storage",
		Example: heredoc.Doc(`
			# Update external storage configuration
			$ stream-cli chat update-external-storage --name my-s3 --properties '{"storage_type":"s3","bucket":"new-bucket"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}

			_, err = h.DoRequest(cmd.Context(), "PUT", "external_storage/"+name, body)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully updated external storage [%s]\n", name)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Storage name")
	fl.StringP("properties", "p", "", "[required] JSON properties")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-external-storage [name]",
		Short: "Delete an external storage",
		Example: heredoc.Doc(`
			# Delete external storage
			$ stream-cli chat delete-external-storage my-s3
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			_, err = h.DoRequest(cmd.Context(), "DELETE", "external_storage/"+args[0], nil)
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
		Use:   "check-external-storage [name] --output-format [json|tree]",
		Short: "Check an external storage configuration",
		Example: heredoc.Doc(`
			# Check external storage
			$ stream-cli chat check-external-storage my-s3
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			resp, err := h.DoRequest(cmd.Context(), "GET", "external_storage/"+args[0]+"/check", nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}
