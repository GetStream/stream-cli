package stats

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		queryCallStatsCmd(),
		queryAggregateCallStatsCmd(),
		getCallReportCmd(),
		getCallStatsMapCmd(),
		queryCallParticipantsCmd(),
		queryCallParticipantSessionsCmd(),
		getCallParticipantSessionMetricsCmd(),
		getSessionParticipantStatsDetailsCmd(),
		querySessionParticipantStatsCmd(),
		getSessionParticipantStatsTimelineCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func callPath(t, id string) string { return "video/call/" + t + "/" + id }

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

func queryCallStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-call-stats --properties [json]",
		Short: "Query call statistics",
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
			return doJSON(cmd, "POST", "video/call/stats", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[optional] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func queryAggregateCallStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-aggregate-call-stats --properties [json]",
		Short: "Query aggregate call statistics",
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
			return doJSON(cmd, "POST", "video/stats", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[optional] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func getCallReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-call-report --type [type] --id [id]",
		Short: "Get call report",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			return doJSON(cmd, "GET", callPath(t, id)+"/report", nil)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func getCallStatsMapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-call-stats-map --call-type [type] --call-id [id] --session [session]",
		Short: "Map call participants by location",
		RunE: func(cmd *cobra.Command, args []string) error {
			ct, _ := cmd.Flags().GetString("call-type")
			cid, _ := cmd.Flags().GetString("call-id")
			session, _ := cmd.Flags().GetString("session")
			return doJSON(cmd, "GET", "video/call_stats/"+ct+"/"+cid+"/"+session+"/map", nil)
		},
	}
	fl := cmd.Flags()
	fl.String("call-type", "", "[required] Call type")
	fl.String("call-id", "", "[required] Call ID")
	fl.StringP("session", "s", "", "[required] Session ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("call-type")
	_ = cmd.MarkFlagRequired("call-id")
	_ = cmd.MarkFlagRequired("session")
	return cmd
}

func queryCallParticipantsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-call-participants --type [type] --id [id] --properties [json]",
		Short: "Query call participants",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
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
			return doJSON(cmd, "POST", callPath(t, id)+"/participants", body)
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("properties", "p", "", "[optional] Query as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func queryCallParticipantSessionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-call-participant-sessions --type [type] --id [id] --session [session]",
		Short: "Query call participant sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			session, _ := cmd.Flags().GetString("session")
			return doJSON(cmd, "POST", callPath(t, id)+"/session/"+session+"/participant_sessions", map[string]interface{}{})
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("session", "s", "", "[required] Session ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("session")
	return cmd
}

func getCallParticipantSessionMetricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-call-participant-session-metrics --type [type] --id [id] --session [session]",
		Short: "Get call participant session metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			session, _ := cmd.Flags().GetString("session")
			return doJSON(cmd, "POST", callPath(t, id)+"/session/"+session+"/participant_sessions", map[string]interface{}{})
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("session", "s", "", "[required] Session ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("session")
	return cmd
}

func getSessionParticipantStatsDetailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-session-participant-stats-details --call-type [type] --call-id [id] --session [session] --user [user] --user-session [user-session]",
		Short: "Get call session participant stats details",
		RunE: func(cmd *cobra.Command, args []string) error {
			ct, _ := cmd.Flags().GetString("call-type")
			cid, _ := cmd.Flags().GetString("call-id")
			session, _ := cmd.Flags().GetString("session")
			user, _ := cmd.Flags().GetString("user")
			userSession, _ := cmd.Flags().GetString("user-session")
			path := "video/call_stats/" + ct + "/" + cid + "/" + session + "/participant/" + user + "/" + userSession + "/details"
			return doJSON(cmd, "GET", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.String("call-type", "", "[required] Call type")
	fl.String("call-id", "", "[required] Call ID")
	fl.StringP("session", "s", "", "[required] Session ID")
	fl.String("user", "", "[required] User ID")
	fl.String("user-session", "", "[required] User session ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("call-type")
	_ = cmd.MarkFlagRequired("call-id")
	_ = cmd.MarkFlagRequired("session")
	_ = cmd.MarkFlagRequired("user")
	_ = cmd.MarkFlagRequired("user-session")
	return cmd
}

func querySessionParticipantStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-session-participant-stats --call-type [type] --call-id [id] --session [session]",
		Short: "Query call participant statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			ct, _ := cmd.Flags().GetString("call-type")
			cid, _ := cmd.Flags().GetString("call-id")
			session, _ := cmd.Flags().GetString("session")
			path := "video/call_stats/" + ct + "/" + cid + "/" + session + "/participants"
			return doJSON(cmd, "GET", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.String("call-type", "", "[required] Call type")
	fl.String("call-id", "", "[required] Call ID")
	fl.StringP("session", "s", "", "[required] Session ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("call-type")
	_ = cmd.MarkFlagRequired("call-id")
	_ = cmd.MarkFlagRequired("session")
	return cmd
}

func getSessionParticipantStatsTimelineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-session-participant-stats-timeline --call-type [type] --call-id [id] --session [session] --user [user] --user-session [user-session]",
		Short: "Get participant timeline events",
		RunE: func(cmd *cobra.Command, args []string) error {
			ct, _ := cmd.Flags().GetString("call-type")
			cid, _ := cmd.Flags().GetString("call-id")
			session, _ := cmd.Flags().GetString("session")
			user, _ := cmd.Flags().GetString("user")
			userSession, _ := cmd.Flags().GetString("user-session")
			path := "video/call_stats/" + ct + "/" + cid + "/" + session + "/participants/" + user + "/" + userSession + "/timeline"
			return doJSON(cmd, "GET", path, nil)
		},
	}
	fl := cmd.Flags()
	fl.String("call-type", "", "[required] Call type")
	fl.String("call-id", "", "[required] Call ID")
	fl.StringP("session", "s", "", "[required] Session ID")
	fl.String("user", "", "[required] User ID")
	fl.String("user-session", "", "[required] User session ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("call-type")
	_ = cmd.MarkFlagRequired("call-id")
	_ = cmd.MarkFlagRequired("session")
	_ = cmd.MarkFlagRequired("user")
	_ = cmd.MarkFlagRequired("user-session")
	return cmd
}
