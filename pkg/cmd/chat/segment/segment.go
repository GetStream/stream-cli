package segment

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		querySegmentsCmd(),
		getSegmentCmd(),
		deleteSegmentCmd(),
		querySegmentTargetsCmd(),
		deleteSegmentTargetsCmd(),
		segmentTargetExistsCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func querySegmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-segments --filter [raw-json]",
		Short: "Query segments",
		Example: heredoc.Doc(`
			# Query all segments
			$ stream-cli chat query-segments --filter '{}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			filterStr, _ := cmd.Flags().GetString("filter")
			var filter interface{}
			if err := json.Unmarshal([]byte(filterStr), &filter); err != nil {
				return err
			}

			body := map[string]interface{}{"filter": filter}

			resp, err := h.DoRequest(cmd.Context(), "POST", "segments/query", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("filter", "f", "{}", "[required] Filter conditions as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func getSegmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-segment [segment-id]",
		Short: "Get a segment by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			resp, err := h.DoRequest(cmd.Context(), "GET", "segments/"+args[0], nil)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func deleteSegmentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delete-segment [segment-id]",
		Short: "Delete a segment",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			_, err = h.DoRequest(cmd.Context(), "DELETE", "segments/"+args[0], nil)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully deleted segment [%s]\n", args[0])
			return nil
		},
	}

	return cmd
}

func querySegmentTargetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-segment-targets --segment-id [id]",
		Short: "Query targets of a segment",
		Example: heredoc.Doc(`
			# Query targets of a segment
			$ stream-cli chat query-segment-targets --segment-id seg-123
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			segID, _ := cmd.Flags().GetString("segment-id")

			resp, err := h.DoRequest(cmd.Context(), "POST", "segments/"+segID+"/targets/query", map[string]interface{}{})
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.String("segment-id", "", "[required] Segment ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("segment-id")

	return cmd
}

func deleteSegmentTargetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-segment-targets --segment-id [id] [target-id-1] [target-id-2] ...",
		Short: "Delete targets from a segment",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			segID, _ := cmd.Flags().GetString("segment-id")
			body := map[string]interface{}{
				"target_ids": args,
			}

			_, err = h.DoRequest(cmd.Context(), "POST", "segments/"+segID+"/deletetargets", body)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted segment targets")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.String("segment-id", "", "[required] Segment ID")
	_ = cmd.MarkFlagRequired("segment-id")

	return cmd
}

func segmentTargetExistsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "segment-target-exists --segment-id [id] --target-id [target-id]",
		Short: "Check if a target exists in a segment",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			segID, _ := cmd.Flags().GetString("segment-id")
			targetID, _ := cmd.Flags().GetString("target-id")

			_, err = h.DoRequest(cmd.Context(), "GET", "segments/"+segID+"/target/"+targetID, nil)
			if err != nil {
				return err
			}

			cmd.Println("Target exists in segment")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.String("segment-id", "", "[required] Segment ID")
	fl.String("target-id", "", "[required] Target ID")
	_ = cmd.MarkFlagRequired("segment-id")
	_ = cmd.MarkFlagRequired("target-id")

	return cmd
}
