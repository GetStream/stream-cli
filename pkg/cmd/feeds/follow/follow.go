package follow

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
		followCmd(),
		unfollowCmd(),
		updateFollowCmd(),
		followBatchCmd(),
		upsertFollowsCmd(),
		queryFollowsCmd(),
		acceptFollowCmd(),
		rejectFollowCmd(),
		unfollowBatchCmd(),
		upsertUnfollowsCmd(),
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

func followCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "follow --properties [json]",
		Short: "Create a follow relationship",
		Example: heredoc.Doc(`
			$ stream-cli feeds follow --properties '{"source":"timeline:user-1","target":"user:user-2"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/follows", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Follow properties as JSON (source, target)")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func unfollowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfollow --source [source] --target [target]",
		Short: "Remove a follow relationship",
		Example: heredoc.Doc(`
			$ stream-cli feeds unfollow --source timeline:user-1 --target user:user-2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			source, _ := cmd.Flags().GetString("source")
			target, _ := cmd.Flags().GetString("target")
			return doJSON(cmd, "DELETE", fmt.Sprintf("api/v2/feeds/follows/%s/%s", source, target), nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("source", "s", "", "[required] Source feed (e.g. timeline:user-1)")
	fl.StringP("target", "t", "", "[required] Target feed (e.g. user:user-2)")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func updateFollowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-follow --properties [json]",
		Short: "Update a follow relationship",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "PATCH", "api/v2/feeds/follows", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Follow update as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func followBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "follow-batch --properties [json]",
		Short: "Create multiple follows at once",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/follows/batch", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Batch follows as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func upsertFollowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-follows --properties [json]",
		Short: "Create or update multiple follows (idempotent)",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/follows/batch/upsert", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Batch follows as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func queryFollowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-follows --properties [json]",
		Short: "Query follow relationships",
		Example: heredoc.Doc(`
			$ stream-cli feeds query-follows --properties '{"filter":{"source_feed":"timeline:user-1"}}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/follows/query", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func acceptFollowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-follow --properties [json]",
		Short: "Accept a pending follow request",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/follows/accept", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Accept request as JSON (source, target)")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func rejectFollowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject-follow --properties [json]",
		Short: "Reject a pending follow request",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/follows/reject", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Reject request as JSON (source, target)")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func unfollowBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfollow-batch --properties [json]",
		Short: "Remove multiple follows at once",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/unfollow/batch", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Batch unfollow as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func upsertUnfollowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-unfollows --properties [json]",
		Short: "Remove multiple follows (idempotent)",
		RunE: func(cmd *cobra.Command, args []string) error {
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			return doJSON(cmd, "POST", "api/v2/feeds/unfollow/batch/upsert", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Batch unfollow as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}
