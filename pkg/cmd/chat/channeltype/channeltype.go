package channeltype

import (
	"encoding/json"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		getCmd(),
		createCmd(),
		deleteCmd(),
		updateCmd(),
		listCmd(),
	}
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-channel-type [channel-type] --output-format [json|tree]",
		Short: "Get channel type",
		Example: heredoc.Doc(`
			# Returns a channel type and prints it as JSON
			$ stream-cli chat get-channel-type livestream

			# Returns a channel type and prints it as a browsable tree
			$ stream-cli chat get-channel-type messaging --output-format tree
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			ct := args[0]

			resp, err := c.GetChannelType(cmd.Context(), ct)
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-channel-type --properties [raw-json-properties]",
		Short: "Create channel type",
		Long: heredoc.Doc(`
			This command creates a new channel type. The 'properties' are raw JSON string.
			The available properties can be found in the Go SDK's 'ChannelType' struct.
		`),
		Example: heredoc.Doc(`
			# Create a new channel type called my-ch-type
			$ stream-cli chat create-channel-type -p "{\"name\": \"my-ch-type\"}"

			# Create a new channel type called reactionless with reactions disabled
			$ stream-cli chat create-channel-type -p "{\"name\": \"reactionless\", \"reactions\": false}"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			properties, _ := cmd.Flags().GetString("properties")

			ct := &stream.ChannelType{ChannelConfig: stream.DefaultChannelConfig}
			err = json.Unmarshal([]byte(properties), ct)
			if err != nil {
				return err
			}

			_, err = c.CreateChannelType(cmd.Context(), ct)
			if err != nil {
				return err
			}

			cmd.Println("Successfully created channel type.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Raw JSON properties")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-channel-type [channel-type]",
		Short: "Delete channel type",
		Long: heredoc.Doc(`
			This command deletes a channel type.
		`),
		Example: heredoc.Doc(`
			# Delete a channel type called my-ch-type
			$ stream-cli chat delete-channel-type my-ch-type
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			ct := args[0]

			_, err = c.GetChannelType(cmd.Context(), ct)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted channel type.")
			return nil
		},
	}
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-channel-type --type [channel-type] --properties [raw-json-properties]",
		Short: "Update channel type",
		Long: heredoc.Doc(`
			This command updates an existing channel type. The 'properties' are raw JSON string.
			The available fields can be checked here:
			https://getstream.io/chat/docs/rest/#channel-types-updatechanneltype
		`),
		Example: heredoc.Doc(`
			# Enabling quotes in an existing channel type
			$ stream-cli chat update-channel-type --type my-channel-type --properties '{"quotes": true}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			ct, _ := cmd.Flags().GetString("type")
			properties, _ := cmd.Flags().GetString("properties")

			props := make(map[string]interface{})
			err = json.Unmarshal([]byte(properties), &props)
			if err != nil {
				return err
			}

			_, err = c.UpdateChannelType(cmd.Context(), ct, props)
			if err != nil {
				return err
			}

			cmd.Println("Successfully updated channel type.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type")
	fl.StringP("properties", "p", "", "[required] Raw JSON properties")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("properties")

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-channel-types --output-format [json|tree]",
		Short: "List channel types",
		Long: heredoc.Doc(`
			This command lists all channel types, including built-in and custom ones.
		`),
		Example: heredoc.Doc(`
			# List all channel types as json (default)
			$ stream-cli chat list-channel-types

			# List all channel types as browsable tree
			$ stream-cli chat list-channel-types --output-format tree
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.ListChannelTypes(cmd.Context())
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp.ChannelTypes)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}
