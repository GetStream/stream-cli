package broadcast

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		startHLSCmd(),
		stopHLSCmd(),
		startRTMPCmd(),
		stopRTMPCmd(),
		stopAllRTMPCmd(),
		startFrameRecordingCmd(),
		stopFrameRecordingCmd(),
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

func startHLSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-hls-broadcasting --type [type] --id [id]",
		Short: "Start HLS broadcasting",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			resp, err := h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/start_broadcasting", map[string]interface{}{})
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

func stopHLSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-hls-broadcasting --type [type] --id [id]",
		Short: "Stop HLS broadcasting",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			_, err = h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/stop_broadcasting", map[string]interface{}{})
			if err != nil {
				return err
			}
			cmd.Println("Successfully stopped HLS broadcasting")
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

func startRTMPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-rtmp-broadcasts --type [type] --id [id] --properties [json]",
		Short: "Start RTMP broadcasts",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			resp, err := h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/rtmp_broadcasts", body)
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
	fl.StringP("properties", "p", "", "[required] RTMP broadcast config as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func stopRTMPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-rtmp-broadcast --type [type] --id [id] --name [name]",
		Short: "Stop a specific RTMP broadcast",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			_, err = h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/rtmp_broadcasts/"+name+"/stop", map[string]interface{}{})
			if err != nil {
				return err
			}
			cmd.Println("Successfully stopped RTMP broadcast")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Call type")
	fl.StringP("id", "i", "", "[required] Call ID")
	fl.StringP("name", "n", "", "[required] Broadcast name")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func stopAllRTMPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-all-rtmp-broadcasts --type [type] --id [id]",
		Short: "Stop all RTMP broadcasts for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			_, err = h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/rtmp_broadcasts/stop", map[string]interface{}{})
			if err != nil {
				return err
			}
			cmd.Println("Successfully stopped all RTMP broadcasts")
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

func startFrameRecordingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-frame-recording --type [type] --id [id]",
		Short: "Start frame recording for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			resp, err := h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/start_frame_recording", map[string]interface{}{})
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

func stopFrameRecordingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-frame-recording --type [type] --id [id]",
		Short: "Stop frame recording for a call",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			t, _ := cmd.Flags().GetString("type")
			id, _ := cmd.Flags().GetString("id")
			_, err = h.DoRequest(cmd.Context(), "POST", callPath(t, id)+"/stop_frame_recording", map[string]interface{}{})
			if err != nil {
				return err
			}
			cmd.Println("Successfully stopped frame recording")
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
