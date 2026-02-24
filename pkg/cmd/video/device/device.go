package device

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listDevicesCmd(),
		createDeviceCmd(),
		deleteDeviceCmd(),
	}
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func listDevicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-devices --user-id [user-id]",
		Short: "List devices for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			userID, _ := cmd.Flags().GetString("user-id")
			path := "video/devices"
			if userID != "" {
				path += "?user_id=" + userID
			}
			resp, err := h.DoRequest(cmd.Context(), "GET", path, nil)
			if err != nil {
				return err
			}
			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}
	fl := cmd.Flags()
	fl.StringP("user-id", "u", "", "[optional] User ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	return cmd
}

func createDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-device --properties [json]",
		Short: "Create a device",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			props, _ := cmd.Flags().GetString("properties")
			var body interface{}
			if err := json.Unmarshal([]byte(props), &body); err != nil {
				return err
			}
			_, err = h.DoRequest(cmd.Context(), "POST", "video/devices", body)
			if err != nil {
				return err
			}
			cmd.Println("Successfully created device")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Device properties as JSON")
	_ = cmd.MarkFlagRequired("properties")
	return cmd
}

func deleteDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-device --device-id [id] --user-id [user-id]",
		Short: "Delete a device",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}
			deviceID, _ := cmd.Flags().GetString("device-id")
			userID, _ := cmd.Flags().GetString("user-id")
			path := "video/devices?id=" + deviceID
			if userID != "" {
				path += "&user_id=" + userID
			}
			_, err = h.DoRequest(cmd.Context(), "DELETE", path, nil)
			if err != nil {
				return err
			}
			cmd.Println("Successfully deleted device")
			return nil
		},
	}
	fl := cmd.Flags()
	fl.String("device-id", "", "[required] Device ID")
	fl.StringP("user-id", "u", "", "[optional] User ID")
	_ = cmd.MarkFlagRequired("device-id")
	return cmd
}
