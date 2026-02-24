package transcription

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listTranscriptionsCmd(),
		startTranscriptionCmd(),
		stopTranscriptionCmd(),
		deleteTranscriptionCmd(),
		startClosedCaptionsCmd(),
		stopClosedCaptionsCmd(),
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

func listTranscriptionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-transcriptions --type [type] --id [id]",
		Short: "List transcriptions for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			resp, err := h.DoRequest(cmd.Context(), "GET", callPath(t, id)+"/transcriptions", nil)
			if err != nil {
				return err
			}
			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
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

func startTranscriptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-transcription --type [type] --id [id]",
		Short: "Start transcription for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			resp, err := h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/start_transcription", map[string]interface{}{})
			if err != nil {
				return err
			}
			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
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

func stopTranscriptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-transcription --type [type] --id [id]",
		Short: "Stop transcription for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			_, err = h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/stop_transcription", map[string]interface{}{})
			if err != nil {
				return err
			}
			cmd.Println("Successfully stopped transcription")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func deleteTranscriptionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-transcription --type [type] --id [id] --session [session] --filename [filename]",
		Short: "Delete a transcription",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			session, _ := cmd.Flags().GetString("session")
			filename, _ := cmd.Flags().GetString("filename")
			_, err = h.DoRequest(cmd.Context(), "DELETE", callPath(t, id)+"/"+session+"/transcriptions/"+filename, nil)
			if err != nil {
				return err
			}
			cmd.Println("Successfully deleted transcription")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("session", "s", "", "[required] Session ID")
	fl.StringP("filename", "f", "", "[required] Transcription filename")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("session")
	_ = cmd.MarkFlagRequired("filename")
	return cmd
}

func startClosedCaptionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-closed-captions --type [type] --id [id]",
		Short: "Start closed captions for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			resp, err := h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/start_closed_captions", map[string]interface{}{})
			if err != nil {
				return err
			}
			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
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

func stopClosedCaptionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-closed-captions --type [type] --id [id]",
		Short: "Stop closed captions for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			_, err = h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/stop_closed_captions", map[string]interface{}{})
			if err != nil {
				return err
			}
			cmd.Println("Successfully stopped closed captions")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
