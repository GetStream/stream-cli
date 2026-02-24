package reminder

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		createReminderCmd(),
		updateReminderCmd(),
		deleteReminderCmd(),
		queryRemindersCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func createReminderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-reminder --message-id [id] --remind-at [RFC3339-timestamp]",
		Short: "Create a reminder for a message",
		Example: heredoc.Doc(`
			# Create a reminder for a message
			$ stream-cli chat create-reminder --message-id msg-1 --remind-at "2025-12-01T10:00:00Z"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			msgID, _ := cmd.Flags().GetString("message-id")
			remindAt, _ := cmd.Flags().GetString("remind-at")

			body := map[string]interface{}{}
			if remindAt != "" {
				body["remind_at"] = remindAt
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "messages/"+msgID+"/reminders", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("message-id", "m", "", "[required] Message ID")
	fl.String("remind-at", "", "[optional] Reminder time in RFC3339 format")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("message-id")

	return cmd
}

func updateReminderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-reminder --message-id [id] --remind-at [RFC3339-timestamp]",
		Short: "Update an existing reminder",
		Example: heredoc.Doc(`
			# Update a reminder
			$ stream-cli chat update-reminder --message-id msg-1 --remind-at "2025-12-15T10:00:00Z"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			msgID, _ := cmd.Flags().GetString("message-id")
			remindAt, _ := cmd.Flags().GetString("remind-at")

			body := map[string]interface{}{}
			if remindAt != "" {
				body["remind_at"] = remindAt
			}

			resp, err := h.DoRequest(cmd.Context(), "PATCH", "messages/"+msgID+"/reminders", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("message-id", "m", "", "[required] Message ID")
	fl.String("remind-at", "", "[optional] Reminder time in RFC3339 format")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("message-id")

	return cmd
}

func deleteReminderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-reminder --message-id [id]",
		Short: "Delete a reminder",
		Example: heredoc.Doc(`
			# Delete a reminder
			$ stream-cli chat delete-reminder --message-id msg-1
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			msgID, _ := cmd.Flags().GetString("message-id")

			_, err = h.DoRequest(cmd.Context(), "DELETE", "messages/"+msgID+"/reminders", nil)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted reminder")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("message-id", "m", "", "[required] Message ID")
	_ = cmd.MarkFlagRequired("message-id")

	return cmd
}

func queryRemindersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-reminders --user-id [user-id]",
		Short: "Query reminders",
		Example: heredoc.Doc(`
			# Query reminders for a user
			$ stream-cli chat query-reminders --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user-id")
			body := map[string]interface{}{
				"user_id": userID,
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "reminders/query", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("user-id", "u", "", "[required] User ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}
