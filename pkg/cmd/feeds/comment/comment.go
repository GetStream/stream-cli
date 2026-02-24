package comment

import (
	"encoding/json"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		addCommentCmd(),
		getCommentCmd(),
		updateCommentCmd(),
		deleteCommentCmd(),
		getCommentsCmd(),
		queryCommentsCmd(),
		addCommentsBatchCmd(),
		getCommentRepliesCmd(),
		addCommentReactionCmd(),
		deleteCommentReactionCmd(),
		queryCommentReactionsCmd(),
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

func addCommentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-comment --properties [json]",
		Short: "Add a comment or reply to an object",
		Example: heredoc.Doc(`
			$ stream-cli feeds add-comment --properties '{"object_type":"activity","object_id":"act-1","comment":"Great post!"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/comments", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Comment properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func getCommentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-comment --id [id]",
		Short: "Get a comment by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "GET", "api/v2/feeds/comments/"+id, nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Comment ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func updateCommentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-comment --id [id] --properties [json]",
		Short: "Update a comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", "api/v2/feeds/comments/"+id, body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Comment ID")
	fl.StringP("properties", "p", "", "[required] Comment update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteCommentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-comment --id [id]",
		Short: "Delete a comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			hard, _ := cmd.Flags().GetBool("hard")
			path := "api/v2/feeds/comments/" + id
			if hard {
				path += "?hard_delete=true"
			}
			return doAction(cmd, "DELETE", path, nil, "Successfully deleted comment")
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Comment ID")
	fl.Bool("hard", false, "[optional] Hard delete")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func getCommentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-comments --object-type [type] --object-id [id]",
		Short: "Get threaded comments for an object",
		Example: heredoc.Doc(`
			$ stream-cli feeds get-comments --object-type activity --object-id act-1
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			objectType, _ := cmd.Flags().GetString("object-type")
			objectID, _ := cmd.Flags().GetString("object-id")
			path := fmt.Sprintf("api/v2/feeds/comments?object_type=%s&object_id=%s", objectType, objectID)
			depth, _ := cmd.Flags().GetInt("depth")
			if depth > 0 {
				path += fmt.Sprintf("&depth=%d", depth)
			}
			sort, _ := cmd.Flags().GetString("sort")
			if sort != "" {
				path += "&sort=" + sort
			}
			limit, _ := cmd.Flags().GetInt("limit")
			if limit > 0 {
				path += fmt.Sprintf("&limit=%d", limit)
			}
			return doJSON(cmd, "GET", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.String("object-type", "", "[required] Object type (e.g. activity)")
	fl.String("object-id", "", "[required] Object ID")
	fl.Int("depth", 0, "[optional] Maximum nested depth")
	fl.String("sort", "", "[optional] Sort order: first, last, top, best, controversial")
	fl.IntP("limit", "l", 0, "[optional] Limit results")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("object-type")
	_ = cmd.MarkFlagRequired("object-id")
	return cmd
}

func queryCommentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-comments --properties [json]",
		Short: "Query comments with filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/comments/query", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func addCommentsBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-comments-batch --properties [json]",
		Short: "Add multiple comments in a batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/comments/batch", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Batch comments as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func getCommentRepliesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-comment-replies --id [id]",
		Short: "Get replies for a comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			path := "api/v2/feeds/comments/" + id + "/replies"
			depth, _ := cmd.Flags().GetInt("depth")
			if depth > 0 {
				path += fmt.Sprintf("?depth=%d", depth)
			}
			return doJSON(cmd, "GET", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Comment ID")
	fl.Int("depth", 0, "[optional] Maximum nested depth")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func addCommentReactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-comment-reaction --id [comment-id] --properties [json]",
		Short: "Add a reaction to a comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/comments/"+id+"/reactions", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Comment ID")
	fl.StringP("properties", "p", "", "[required] Reaction properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteCommentReactionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-comment-reaction --id [comment-id] --type [type]",
		Short: "Remove a reaction from a comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			reactionType, _ := cmd.Flags().GetString("type")
			return doJSON(cmd, "DELETE", fmt.Sprintf("api/v2/feeds/comments/%s/reactions/%s", id, reactionType), nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Comment ID")
	fl.StringP("type", "t", "", "[required] Reaction type")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func queryCommentReactionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-comment-reactions --id [comment-id] --properties [json]",
		Short: "Query reactions on a comment",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if props != "" {
				if err := json.Unmarshal([]byte(props), &body); err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}
			return doJSON(cmd, "POST", "api/v2/feeds/comments/"+id+"/reactions/query", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Comment ID")
	fl.StringP("properties", "p", "", "[optional] Query filters as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
