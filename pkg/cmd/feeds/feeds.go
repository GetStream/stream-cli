package feeds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

// ImportData represents the structure of feeds import data
type ImportData struct {
	Users      []User     `json:"users,omitempty"`
	Feeds      []Feed     `json:"feeds,omitempty"`
	Activities []Activity `json:"activities,omitempty"`
	Follows    []User     `json:"follows,omitempty"`
	Members    []Member   `json:"members,omitempty"`
	Comments   []User     `json:"comments,omitempty"`
	Reactions  []Reaction `json:"reactions,omitempty"`
}

// User represents a user in the import data
type User struct {
	ID     string         `json:"id"`
	Role   string         `json:"role,omitempty"`
	Custom map[string]any `json:"custom,omitempty"`
}

// Feed represents a feed in the import data
type Feed struct {
	GroupID     string                 `json:"group_id"`
	ID          string                 `json:"id"`
	FID         string                 `json:"fid"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
	FilterTags  []string               `json:"filter_tags,omitempty"`
	Visibility  string                 `json:"visibility,omitempty"`
	CreatedAt   int64                  `json:"created_at,omitempty"`
	UpdatedAt   int64                  `json:"updated_at,omitempty"`
	DeletedAt   *int64                 `json:"deleted_at,omitempty"`
}

// Activity represents an activity in the import data
type Activity struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	Feeds          []string               `json:"feeds"`
	Visibility     string                 `json:"visibility,omitempty"`
	CreatedAt      int64                  `json:"created_at,omitempty"`
	UpdatedAt      int64                  `json:"updated_at,omitempty"`
	Attachments    []interface{}          `json:"attachments,omitempty"`
	MentionedUsers []string               `json:"mentioned_users,omitempty"`
	Custom         map[string]interface{} `json:"custom,omitempty"`
	Text           string                 `json:"text,omitempty"`
	SearchData     map[string]interface{} `json:"search_data,omitempty"`
	FilterTags     []string               `json:"filter_tags,omitempty"`
	InterestTags   []string               `json:"interest_tags,omitempty"`
}

// Member represents a member in the import data
type Member struct {
	UserID string         `json:"user_id"`
	Feed   string         `json:"feed"`
	Custom map[string]any `json:"custom,omitempty"`
	Role   string         `json:"role,omitempty"`
}

// Reaction represents a reaction in the import data
type Reaction struct {
	ActivityID string         `json:"activity_id"`
	CommentID  string         `json:"comment_id,omitempty"`
	Type       string         `json:"type"`
	UserID     string         `json:"user_id"`
	Custom     map[string]any `json:"custom,omitempty"`
}

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		importValidateCmd(),
		importCmd(),
		importStatusCmd(),
	}
}

func validateFeedsFile(cmd *cobra.Command, filename string) error {
	reader, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer reader.Close()

	// Parse JSON and validate structure
	var data ImportData
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate that at least one section is present
	if len(data.Users) == 0 && len(data.Feeds) == 0 && len(data.Activities) == 0 &&
		len(data.Follows) == 0 && len(data.Members) == 0 && len(data.Comments) == 0 &&
		len(data.Reactions) == 0 {
		return fmt.Errorf("file must contain at least one of: users, feeds, activities, follows, members, comments, or reactions")
	}

	// Validate users
	if err := validateUsers(data.Users); err != nil {
		return fmt.Errorf("users validation failed: %w", err)
	}

	// Validate feeds
	if err := validateFeeds(data.Feeds); err != nil {
		return fmt.Errorf("feeds validation failed: %w", err)
	}

	// Validate activities
	if err := validateActivities(data.Activities); err != nil {
		return fmt.Errorf("activities validation failed: %w", err)
	}

	// Validate follows
	if err := validateUsers(data.Follows); err != nil {
		return fmt.Errorf("follows validation failed: %w", err)
	}

	// Validate members
	if err := validateMembers(data.Members); err != nil {
		return fmt.Errorf("members validation failed: %w", err)
	}

	// Validate comments
	if err := validateUsers(data.Comments); err != nil {
		return fmt.Errorf("comments validation failed: %w", err)
	}

	// Validate reactions
	if err := validateReactions(data.Reactions); err != nil {
		return fmt.Errorf("reactions validation failed: %w", err)
	}

	// Only print success message if all validations pass
	cmd.Printf("✅ File '%s' is valid JSON with proper feeds import structure\n", filename)
	cmd.Printf("📊 Import summary:\n")
	if len(data.Users) > 0 {
		cmd.Printf("   - Users: %d\n", len(data.Users))
	}
	if len(data.Feeds) > 0 {
		cmd.Printf("   - Feeds: %d\n", len(data.Feeds))
	}
	if len(data.Activities) > 0 {
		cmd.Printf("   - Activities: %d\n", len(data.Activities))
	}
	if len(data.Follows) > 0 {
		cmd.Printf("   - Follows: %d\n", len(data.Follows))
	}
	if len(data.Members) > 0 {
		cmd.Printf("   - Members: %d\n", len(data.Members))
	}
	if len(data.Comments) > 0 {
		cmd.Printf("   - Comments: %d\n", len(data.Comments))
	}
	if len(data.Reactions) > 0 {
		cmd.Printf("   - Reactions: %d\n", len(data.Reactions))
	}

	return nil
}

func validateUsers(users []User) error {
	for i, user := range users {
		if user.ID == "" {
			return fmt.Errorf("user at index %d: id is required", i)
		}
	}
	return nil
}

func validateFeeds(feeds []Feed) error {
	for i, feed := range feeds {
		if feed.ID == "" {
			return fmt.Errorf("feed at index %d: id is required", i)
		}
		if feed.GroupID == "" {
			return fmt.Errorf("feed at index %d: group_id is required", i)
		}
		if feed.FID == "" {
			return fmt.Errorf("feed at index %d: fid is required", i)
		}
		if feed.Name == "" {
			return fmt.Errorf("feed at index %d: name is required", i)
		}
	}
	return nil
}

func validateActivities(activities []Activity) error {
	for i, activity := range activities {
		if activity.ID == "" {
			return fmt.Errorf("activity at index %d: id is required", i)
		}
		if activity.Type == "" {
			return fmt.Errorf("activity at index %d: type is required", i)
		}
		if len(activity.Feeds) == 0 {
			return fmt.Errorf("activity at index %d: feeds is required and cannot be empty", i)
		}
	}
	return nil
}

func validateMembers(members []Member) error {
	for i, member := range members {
		if member.UserID == "" {
			return fmt.Errorf("member at index %d: user_id is required", i)
		}
		if member.Feed == "" {
			return fmt.Errorf("member at index %d: feed is required", i)
		}
	}
	return nil
}

func validateReactions(reactions []Reaction) error {
	for i, reaction := range reactions {
		if reaction.ActivityID == "" {
			return fmt.Errorf("reaction at index %d: activity_id is required", i)
		}
		if reaction.Type == "" {
			return fmt.Errorf("reaction at index %d: type is required", i)
		}
		if reaction.UserID == "" {
			return fmt.Errorf("reaction at index %d: user_id is required", i)
		}
	}
	return nil
}

func importValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-validate [filename]",
		Short: "Validate feeds import file",
		Long: heredoc.Doc(`
			Validates a JSON file for feeds import.
			This command checks if the file is valid JSON format.
		`),
		Example: heredoc.Doc(`
			# Validates a JSON feeds import file
			$ stream-cli feeds import-validate feeds-data.json
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			return validateFeedsFile(cmd, filename)
		},
	}

	return cmd
}

func uploadFeedsToS3(ctx context.Context, filename, url string) error {
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %d", resp.StatusCode)
	}

	return nil
}

func importCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [filename] --apikey [api-key]",
		Short: "Import feeds data",
		Long: heredoc.Doc(`
			Imports feeds data from a JSON file.
			This command uploads the file to S3 and initiates the import process.
		`),
		Example: heredoc.Doc(`
			# Import feeds data with API key
			$ stream-cli feeds import feeds-data.json --apikey your-api-key

			# Import feeds data with custom mode
			$ stream-cli feeds import feeds-data.json --apikey your-api-key --mode insert
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			filename := args[0]

			// Validate file first
			if err := validateFeedsFile(cmd, filename); err != nil {
				return err
			}

			// Create import URL
			createImportURLResp, err := c.CreateImportURL(cmd.Context(), filepath.Base(filename))
			if err != nil {
				return err
			}

			// Upload to S3
			if err := uploadFeedsToS3(cmd.Context(), filename, createImportURLResp.UploadURL); err != nil {
				return err
			}

			// Determine import mode
			mode := stream.UpsertMode
			if m, _ := cmd.Flags().GetString("mode"); stream.ImportMode(m) == stream.InsertMode {
				mode = stream.InsertMode
			}

			// Create import
			createImportResp, err := c.CreateImport(cmd.Context(), createImportURLResp.Path, mode)
			if err != nil {
				return err
			}

			cmd.Printf("✅ Import started successfully\n")
			cmd.Printf("Import ID: %s\n", createImportResp.ImportTask.ID)

			return utils.PrintObject(cmd, createImportResp.ImportTask)
		},
	}

	fl := cmd.Flags()
	fl.StringP("apikey", "k", "", "[required] API key for authentication")
	fl.StringP("mode", "m", "upsert", "[optional] Import mode. Can be upsert or insert")
	_ = cmd.MarkFlagRequired("apikey")

	return cmd
}

func importStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-status [import-id]",
		Short: "Check import status",
		Long: heredoc.Doc(`
			Checks the status of a feeds import operation.
			You can optionally watch for completion with the --watch flag.
		`),
		Example: heredoc.Doc(`
			# Check import status
			$ stream-cli feeds import-status dcb6e366-93ec-4e52-af6f-b0c030ad5272

			# Watch import until completion
			$ stream-cli feeds import-status dcb6e366-93ec-4e52-af6f-b0c030ad5272 --watch
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

				// Wait before checking again
				time.Sleep(5 * time.Second)
			}

			return nil
		},
	}

	fl := cmd.Flags()
	fl.BoolP("watch", "w", false, "[optional] Watch import until completion")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}
