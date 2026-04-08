package imports

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	getstream "github.com/GetStream/getstream-go/v4"
	"github.com/MakeNowJust/heredoc"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		uploadCmd(),
		getCmd(),
		listCmd(),
	}
}

// getImportV2Task works around a server-side issue where GET /api/v2/imports/v2/{id}
// requires a JSON body. The SDK sends no body for GET requests, causing a 400 error.
func getImportV2Task(ctx context.Context, app *config.App, id string) (map[string]any, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"server": true})
	authToken, err := token.SignedString([]byte(app.AccessSecretKey))
	if err != nil {
		return nil, fmt.Errorf("creating auth token: %w", err)
	}

	baseURL := app.ChatURL
	if baseURL == "" {
		baseURL = config.DefaultChatEdgeURL
	}

	reqURL := fmt.Sprintf("%s/api/v2/imports/v2/%s?api_key=%s",
		baseURL, url.PathEscape(id), url.QueryEscape(app.AccessKey))

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Stream-Auth-Type", "jwt")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
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
		Use:   "upload-import [filename] --output-format [json|tree]",
		Short: "Upload an import for Feeds",
		Example: heredoc.Doc(`
			# Uploads a feeds import and prints it as JSON
			$ stream-cli feeds upload-import data.json

			# Uploads a feeds import and prints it as a browsable tree
			$ stream-cli feeds upload-import data.json --output-format tree
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := config.GetConfig(cmd).GetFeedsClient(cmd)
			if err != nil {
				return err
			}

			filename := args[0]

			createURLResp, err := client.CreateImportURL(cmd.Context(), &getstream.CreateImportURLRequest{
				Filename: getstream.PtrTo(filepath.Base(filename)),
			})
			if err != nil {
				return err
			}
			if err := uploadToS3(cmd.Context(), filename, createURLResp.Data.UploadUrl); err != nil {
				return err
			}

			bucket, region, err := utils.S3BucketAndRegionFromUploadURL(createURLResp.Data.UploadUrl)
			if err != nil {
				return err
			}
			dir := createURLResp.Data.Path

			skipReferencesCheck, err := cmd.Flags().GetBool("skip-references-check")
			if err != nil {
				return err
			}

			resp, err := client.CreateImportV2Task(cmd.Context(), &getstream.CreateImportV2TaskRequest{
				Product: "feeds",
				Settings: getstream.ImportV2TaskSettings{
					SkipReferencesCheck: getstream.PtrTo(skipReferencesCheck),
					S3: &getstream.ImportV2TaskSettingsS3{
						Bucket: getstream.PtrTo(bucket),
						Dir:    &dir,
						Region: getstream.PtrTo(region),
					},
				},
			})
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp.Data)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	fl.Bool("skip-references-check", false, "[optional] Skip references validation for the import (default false)")

	return cmd
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-import [task-id] --output-format [json|tree] --watch",
		Short: "Get a feeds import task",
		Example: heredoc.Doc(`
			# Returns a feeds import and prints it as JSON
			$ stream-cli feeds get-import dcb6e366-93ec-4e52-af6f-b0c030ad5272

			# Returns a feeds import and watches for completion
			$ stream-cli feeds get-import dcb6e366-93ec-4e52-af6f-b0c030ad5272 --watch
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig(cmd)
			app, err := cfg.GetDefaultAppOrExplicit(cmd)
			if err != nil {
				return err
			}

			id := args[0]
			watch, _ := cmd.Flags().GetBool("watch")

			for {
				result, err := getImportV2Task(cmd.Context(), app, id)
				if err != nil {
					return err
				}

				err = utils.PrintObject(cmd, result)
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
		Use:   "list-imports --output-format [json|tree] --state [1-4]",
		Short: "List feeds import tasks",
		Example: heredoc.Doc(`
			# List all feeds imports as json (default)
			$ stream-cli feeds list-imports

			# List feeds imports filtered by state
			$ stream-cli feeds list-imports --state 2

			# List all feeds imports as browsable tree
			$ stream-cli feeds list-imports --output-format tree
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := config.GetConfig(cmd).GetFeedsClient(cmd)
			if err != nil {
				return err
			}

			state, _ := cmd.Flags().GetInt("state")

			req := &getstream.ListImportV2TasksRequest{}
			if state > 0 {
				req.State = getstream.PtrTo(state)
			}

			resp, err := client.ListImportV2Tasks(cmd.Context(), req)
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp.Data)
		},
	}

	fl := cmd.Flags()
	fl.IntP("state", "s", 0, "[optional] Filter imports by state (1-4)")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}
