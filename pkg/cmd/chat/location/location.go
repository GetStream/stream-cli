package location

import (
	"encoding/json"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getLiveLocationsCmd(),
		updateLiveLocationCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func getLiveLocationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-live-locations --user-id [user-id]",
		Short: "Get active live locations for a user",
		Example: heredoc.Doc(`
			# Get live locations for user 'joe'
			$ stream-cli chat get-live-locations --user-id joe
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user-id")

			resp, err := h.DoRequest(cmd.Context(), "GET", "users/live_locations?user_id="+userID, nil)
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

func updateLiveLocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-live-location --message-id [id] --latitude [lat] --longitude [lon]",
		Short: "Update a live location",
		Example: heredoc.Doc(`
			# Update a live location
			$ stream-cli chat update-live-location --message-id msg-1 --latitude 40.7128 --longitude -74.0060
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			msgID, _ := cmd.Flags().GetString("message-id")
			lat, _ := cmd.Flags().GetFloat64("latitude")
			lon, _ := cmd.Flags().GetFloat64("longitude")

			body := map[string]interface{}{
				"message_id": msgID,
				"latitude":   lat,
				"longitude":  lon,
			}

			resp, err := h.DoRequest(cmd.Context(), "PUT", "users/live_locations", body)
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
	fl.Float64("latitude", 0, "[required] Latitude coordinate")
	fl.Float64("longitude", 0, "[required] Longitude coordinate")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("message-id")
	_ = cmd.MarkFlagRequired("latitude")
	_ = cmd.MarkFlagRequired("longitude")

	return cmd
}
