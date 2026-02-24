package recording

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listRecordingsCmd(),
		startRecordingCmd(),
		stopRecordingCmd(),
		deleteRecordingCmd(),
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

func listRecordingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-recordings --type [type] --id [id]",
		Short: "List recordings for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			resp, err := h.DoRequest(cmd.Context(), "GET", callPath(t, id)+"/recordings", nil)
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

func startRecordingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-recording --type [type] --id [id] --recording-type [recording-type]",
		Short: "Start recording a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			recType, _ := cmd.Flags().GetString("recording-type")
			if recType == "" {
				recType = "audio_and_video"
			}
			resp, err := h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/recordings/"+recType+"/start", map[string]interface{}{})
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
	fl.String("recording-type", "audio_and_video", "[optional] Recording type")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func stopRecordingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-recording --type [type] --id [id] --recording-type [recording-type]",
		Short: "Stop recording a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			recType, _ := cmd.Flags().GetString("recording-type")
			if recType == "" {
				recType = "audio_and_video"
			}
			_, err = h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/recordings/"+recType+"/stop", map[string]interface{}{})
			if err != nil {
				return err
			}
			cmd.Println("Successfully stopped recording")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.String("recording-type", "audio_and_video", "[optional] Recording type")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func deleteRecordingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-recording --type [type] --id [id] --session [session] --filename [filename]",
		Short: "Delete a recording",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			session, _ := cmd.Flags().GetString("session")
			filename, _ := cmd.Flags().GetString("filename")
			_, err = h.DoRequest(cmd.Context(), "DELETE", callPath(t, id)+"/"+session+"/recordings/"+filename, nil)
			if err != nil {
				return err
			}
			cmd.Println("Successfully deleted recording")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("session", "s", "", "[required] Session ID")
	fl.StringP("filename", "f", "", "[required] Recording filename")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("session")
	_ = cmd.MarkFlagRequired("filename")
	return cmd
}
