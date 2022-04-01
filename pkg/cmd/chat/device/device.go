package device

import (
	stream_chat "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		createCmd(),
		listCmd(),
		deleteCmd()}
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-device --id [device-id] --push-provider [firebase|apn|xiaomi|huawei] --push-provider-name [provider-name] --user-id [user-id]",
		Short: "Create a device",
		Long: heredoc.Doc(`
			Registering a device associates it with a user and tells
			the push provider to send new message notifications to the device.
		`),
		Example: heredoc.Doc(`
			# Create a device with a firebase push provider
			$ stream-cli chat create-device --id "my-device-id" --push-provider firebase --push-provider-name "my-firebase-project-id" --user-id "my-user-id"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			id, _ := cmd.Flags().GetString("id")
			provider, _ := cmd.Flags().GetString("push-provider")
			providerName, _ := cmd.Flags().GetString("push-provider-name")
			userID, _ := cmd.Flags().GetString("user-id")

			d := &stream_chat.Device{ID: id, PushProvider: provider, PushProviderName: providerName, UserID: userID}

			_, err = c.AddDevice(cmd.Context(), d)
			if err != nil {
				return err
			}

			cmd.Println("Successfully created device")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Device ID")
	fl.StringP("push-provider", "p", "", "[required] Push provider. Can be apn, firebase, xiaomi, huawei")
	fl.StringP("push-provider-name", "n", "", "[optional] Push provider name")
	fl.StringP("user-id", "u", "", "[required] User ID")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("push-provider")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-devices --user-id [user-id] --output-format [json|tree]",
		Short: "List devices",
		Long:  "Provides a list of all devices associated with a user.",
		Example: heredoc.Doc(`
			# List devices for a user
			$ stream-cli chat list-devices --user-id "my-user-id"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user-id")

			devices, err := c.GetDevices(cmd.Context(), userID)
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, devices.Devices)
		},
	}

	fl := cmd.Flags()
	fl.StringP("user-id", "u", "", "[required] User ID")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-device --id [device-id] --user-id [user-id]",
		Short: "Delete a device",
		Long: heredoc.Doc(`
			Unregistering a device removes the device from the user
			and stops further new message notifications.
		`),
		Example: heredoc.Doc(`
			# Delete "my-device-id" device
			$ stream-cli chat delete-device --id "my-device-id" --user-id "my-user-id"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			id, _ := cmd.Flags().GetString("id")
			userID, _ := cmd.Flags().GetString("user-id")

			_, err = c.DeleteDevice(cmd.Context(), userID, id)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted device")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("id", "i", "", "[required] Device ID to delete")
	fl.StringP("user-id", "u", "", "[required] ID of the user who deletes the device")
	_ = cmd.MarkFlagRequired("id")
	_ = cmd.MarkFlagRequired("user-id")

	return cmd
}
