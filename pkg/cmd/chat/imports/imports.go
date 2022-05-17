package imports

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/cmd/chat/imports/validator"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		validateCmd(),
		uploadCmd(),
		getCmd(),
		listCmd(),
	}
}

func validateFile(ctx context.Context, c *stream.Client, filename string) (*validator.Results, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	rolesResp, err := c.Permissions().ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	channelTypesResp, err := c.ListChannelTypes(ctx)
	if err != nil {
		return nil, err
	}

	return validator.New(reader, rolesResp.Roles, channelTypesResp.ChannelTypes).Validate(), nil
}

func validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate-import [filename]",
		Short: "Validate import file",
		Example: heredoc.Doc(`
			# Validates a JSON import file
			$ stream-cli chat validate-import data.json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			results, err := validateFile(cmd.Context(), c, args[0])
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, results)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}

func uploadToS3(ctx context.Context, filename, url string) error {
	data, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer data.Close()

	stat, err := data.Stat()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = stat.Size()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func uploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-import [filename] --mode [upsert|insert] --output-format [json|tree]",
		Short: "Upload an import",
		Example: heredoc.Doc(`
			# Uploads an import and prints it as JSON
			$ stream-cli chat upload-import data.json --mode insert

			# Uploads an import and prints it as a browsable tree
			$ stream-cli chat upload-import data.json --mode insert --output-format tree
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			filename := args[0]

			results, err := validateFile(cmd.Context(), c, filename)
			if err != nil {
				return err
			}
			if results.HasErrors() {
				return utils.PrintObject(cmd, results)
			}

			mode := stream.UpsertMode
			if m, _ := cmd.Flags().GetString("mode"); stream.ImportMode(m) == stream.InsertMode {
				mode = stream.InsertMode
			}

			createImportURLResp, err := c.CreateImportURL(cmd.Context(), filepath.Base(filename))
			if err != nil {
				return err
			}

			if err := uploadToS3(cmd.Context(), filename, createImportURLResp.UploadURL); err != nil {
				return err
			}
			createImportResp, err := c.CreateImport(cmd.Context(), createImportURLResp.Path, mode)
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, createImportResp.ImportTask)
		},
	}

	fl := cmd.Flags()
	fl.StringP("mode", "m", "upsert", "[optional] Import mode. Canbe upsert or insert")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-import [import-id] --output-format [json|tree] --watch",
		Short: "Get import",
		Example: heredoc.Doc(`
			# Returns an import and prints it as JSON
			$ stream-cli chat get-import dcb6e366-93ec-4e52-af6f-b0c030ad5272

			# Returns an import and prints it as JSON, and wait for it to complete
			$ stream-cli chat get-import dcb6e366-93ec-4e52-af6f-b0c030ad5272 --watch
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			id := args[0]
			watch, _ := cmd.Flags().GetBool("watch")

			for {
				resp, err := c.GetImport(cmd.Context(), id)
				if err != nil {
					return err
				}

				err = utils.PrintObject(cmd, resp.ImportTask)
				if err != nil {
					return err
				}

				if !watch {
					break
				}

				time.Sleep(5 * time.Second)
			}

			return nil
		},
	}

	fl := cmd.Flags()
	fl.BoolP("watch", "w", false, "[optional] Keep polling the import to track its status")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-imports --offset [int] --limit [int] --output-format [json|tree]",
		Short: "List imports",
		Example: heredoc.Doc(`
			# List all imports as json (default)
			$ stream-cli chat list-imports

			# List all imports as browsable tree
			$ stream-cli chat list-imports --output-format tree
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			limit, _ := cmd.Flags().GetInt("limit")
			offset, _ := cmd.Flags().GetInt("offset")

			resp, err := c.ListImports(cmd.Context(), &stream.ListImportsOptions{
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp.ImportTasks)
		},
	}

	fl := cmd.Flags()
	fl.IntP("limit", "l", 10, "[optional] The number of imports returned")
	fl.IntP("offset", "O", 0, "[optional] The starting offset of imports returned")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}
