package push

import (
	"encoding/json"

	stream_chat "github.com/GetStream/stream-chat-go/v5"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		upsertCmd(),
		deleteCmd(),
		listCmd(),
		testCmd(),
	}
}

func upsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-pushprovider --properties [raw-json]",
		Short: "Create or updates a push provider",
		Long: `
			The "--properties" parameter expects a raw json string that can be
			unmarshalled into a stream_chat.PushProvider object on the Go SDK side.
			See the example section.
			Available properties:
			type
			name
			description
			disabled_at
			disabled_reason

			apn_auth_key
			apn_key_id
			apn_team_id
			apn_topic

			firebase_notification_template
			firebase_apn_template
			firebase_credentials

			huawei_app_id
			huawei_app_secret

			xiaomi_package_name
			xiaomi_app_secret
		`,
		Example: heredoc.Doc(`
			# Setting up an APN push provider
			$ stream-cli chat upsert-pushprovider --properties "{'type': 'apn', 'name': 'staging', 'apn_auth_key': 'key', 'apn_key_id': 'id', 'apn_topic': 'topic', 'apn_team_id': 'id'}"

			# Setting up a Firebase push provider
			$ stream-cli chat upsert-pushprovider --properties "{'type': 'firebase', 'name': 'staging', 'firebase_credentials': 'credentials'}"

			# Setting up a Huawei push provider
			$ stream-cli chat upsert-pushprovider --properties "{'type': 'huawei', 'name': 'staging', 'huawei_app_id': 'id', 'huawei_app_secret': 'secret'}"

			# Setting up a Xiaomi push provider
			$ stream-cli chat upsert-pushprovider --properties "{'type': 'xiaomi', 'name': 'staging', 'xiaomi_package_name': 'name', 'xiaomi_app_secret': 'secret'}"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			prop, _ := cmd.Flags().GetString("properties")

			var p stream_chat.PushProvider
			err = json.Unmarshal([]byte(prop), &p)
			if err != nil {
				return err
			}

			_, err = c.UpsertPushProvider(cmd.Context(), &p)
			if err != nil {
				return err
			}

			cmd.Println("Successfully updated push provider.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Raw json properties to send to the backend")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-pushproviders --output-format [json|tree]",
		Short: "List all push providers",
		Example: heredoc.Doc(`
			# List all push providers
			$ stream-cli chat list-pushproviders
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.ListPushProviders(cmd.Context())
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp.PushProviders)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "Output format. One of: json|tree")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-pushprovider --push-provider-type [type] --push-provider-name [name]",
		Short: "Delete a push provider",
		Example: heredoc.Doc(`
			# Delete an APN push provider
			$ stream-cli chat delete-pushprovider --push-provider-type apn --push-provider-name staging
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			providerName, _ := cmd.Flags().GetString("push-provider-name")
			providerType, _ := cmd.Flags().GetString("push-provider-type")

			_, err = c.DeletePushProvider(cmd.Context(), providerType, providerName)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted push provider.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("push-provider-type", "t", "", "[required] Push provider type")
	fl.StringP("push-provider-name", "n", "", "[required] Push provider name")
	_ = cmd.MarkFlagRequired("push-provider-type")
	_ = cmd.MarkFlagRequired("push-provider-name")

	return cmd
}

func testCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "test-push --message-id [string]" +
			" --apn-template [string]" +
			" --firebase-template [string]" +
			" --firebase-data-template [string]" +
			" --skip-devices [true|false]" +
			" --push-provider-name [string]" +
			" --push-provider-type [string]" +
			" --user-id [string]" +
			" --output-format [json|tree]",
		Short: "Test push notifications",
		Example: heredoc.Doc(`
			# A test push notification for a certain message id
			$ stream-cli chat test-push --message-id msgid --user-id id --skip-devices true
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			msgID, _ := cmd.Flags().GetString("message-id")
			skipDevices, _ := cmd.Flags().GetBool("skip-devices")
			pushProviderName, _ := cmd.Flags().GetString("push-provider-name")
			pushProviderType, _ := cmd.Flags().GetString("push-provider-type")
			userID, _ := cmd.Flags().GetString("user-id")

			p := &stream_chat.CheckPushRequest{
				MessageID:        msgID,
				SkipDevices:      &skipDevices,
				PushProviderName: pushProviderName,
				PushProviderType: pushProviderType,
				UserID:           userID,
			}

			resp, err := c.CheckPush(cmd.Context(), p)
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp)
		},
	}

	fl := cmd.Flags()
	fl.String("message-id", "", "[optional] Message id to test")
	fl.Bool("skip-devices", false, "[optional] Whether to notify devices")
	fl.String("push-provider-name", "", "[optional] Push provider name to use")
	fl.String("push-provider-type", "", "[optional] Push provider type to use")
	fl.String("user-id", "", "[optional] User id to initiate the test")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}
