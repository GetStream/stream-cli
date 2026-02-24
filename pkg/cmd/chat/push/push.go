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
		checkSNSCmd(),
		checkSQSCmd(),
		getPushTemplatesCmd(),
		upsertPushTemplateCmd(),
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

func checkSNSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-sns --sns-topic-arn [arn] --sns-key [key] --sns-secret [secret]",
		Short: "Validate Amazon SNS configuration",
		Example: heredoc.Doc(`
			# Check SNS configuration
			$ stream-cli chat check-sns --sns-topic-arn arn:aws:sns:us-east-1:123:topic --sns-key key --sns-secret secret
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			arn, _ := cmd.Flags().GetString("sns-topic-arn")
			key, _ := cmd.Flags().GetString("sns-key")
			secret, _ := cmd.Flags().GetString("sns-secret")

			body := map[string]interface{}{
				"sns_topic_arn": arn,
				"sns_key":       key,
				"sns_secret":    secret,
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "check_sns", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.String("sns-topic-arn", "", "[optional] SNS topic ARN")
	fl.String("sns-key", "", "[optional] AWS SNS access key")
	fl.String("sns-secret", "", "[optional] AWS SNS secret key")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func checkSQSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-sqs --sqs-url [url] --sqs-key [key] --sqs-secret [secret]",
		Short: "Validate Amazon SQS configuration",
		Example: heredoc.Doc(`
			# Check SQS configuration
			$ stream-cli chat check-sqs --sqs-url https://sqs.us-east-1.amazonaws.com/123/queue --sqs-key key --sqs-secret secret
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			sqsURL, _ := cmd.Flags().GetString("sqs-url")
			key, _ := cmd.Flags().GetString("sqs-key")
			secret, _ := cmd.Flags().GetString("sqs-secret")

			body := map[string]interface{}{
				"sqs_url":    sqsURL,
				"sqs_key":    key,
				"sqs_secret": secret,
			}

			resp, err := h.DoRequest(cmd.Context(), "POST", "check_sqs", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.String("sqs-url", "", "[optional] SQS endpoint URL")
	fl.String("sqs-key", "", "[optional] AWS SQS access key")
	fl.String("sqs-secret", "", "[optional] AWS SQS secret key")
	fl.StringP("output-format", "o", "json", "[optional] Output format")

	return cmd
}

func getHTTPClient(cmd *cobra.Command) (*utils.HTTPClient, error) {
	appCfg, err := config.GetConfig(cmd).GetAppConfig(cmd)
	if err != nil {
		return nil, err
	}
	return utils.NewHTTPClient(appCfg.AccessKey, appCfg.AccessSecretKey, appCfg.ChatURL), nil
}

func getPushTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-push-templates --push-provider-type [type]",
		Short: "Get push notification templates",
		Example: heredoc.Doc(`
			# Get push templates for Firebase
			$ stream-cli chat get-push-templates --push-provider-type firebase
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := getHTTPClient(cmd)
			if err != nil {
				return err
			}

			providerType, _ := cmd.Flags().GetString("push-provider-type")
			providerName, _ := cmd.Flags().GetString("push-provider-name")

			path := "push_templates?push_provider_type=" + providerType
			if providerName != "" {
				path += "&push_provider_name=" + providerName
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
	fl.String("push-provider-type", "", "[required] Push provider type (firebase, apn, huawei, xiaomi)")
	fl.String("push-provider-name", "", "[optional] Push provider name")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("push-provider-type")

	return cmd
}

func upsertPushTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-push-template --properties [raw-json]",
		Short: "Create or update a push notification template",
		Example: heredoc.Doc(`
			# Upsert a push template
			$ stream-cli chat upsert-push-template --properties '{"event_type":"message.new","push_provider_type":"firebase","template":"{...}"}'
		`),
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

			resp, err := h.DoRequest(cmd.Context(), "POST", "push_templates", body)
			if err != nil {
				return err
			}

			var result interface{}
			_ = json.Unmarshal(resp, &result)
			return utils.PrintObject(cmd, result)
		},
	}

	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Template properties as JSON")
	fl.StringP("output-format", "o", "json", "[optional] Output format")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}
