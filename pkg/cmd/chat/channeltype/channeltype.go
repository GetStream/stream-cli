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
		listCmd()}
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-channel-type --type [channel-type] --output-format [json|tree]",
		Short: "Get channel type",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			ct, _ := cmd.Flags().GetString("type")

			resp, err := c.GetChannelType(cmd.Context(), ct)
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, resp)
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	cmd.MarkFlagRequired("type")

	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-channel-type --properties [raw-json-properties]",
		Short: "Create channel type",
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
	cmd.MarkFlagRequired("properties")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-channel-type --type [channel-type]",
		Short: "Delete channel type",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			ct, _ := cmd.Flags().GetString("type")

			_, err = c.GetChannelType(cmd.Context(), ct)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted channel type.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type")
	cmd.MarkFlagRequired("type")

	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-channel-type --type [channel-type] --properties [raw-json-properties]",
		Short: "Update channel type",
		Example: heredoc.Doc(`
			update-channel-type --type my-channel-type --properties '{"quotes": true}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			ct, _ := cmd.Flags().GetString("type")
			properties, _ := cmd.Flags().GetString("properties")

			props := make(map[string]interface{})
			json.Unmarshal([]byte(properties), &props)

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
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("properties")

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-channel-types --type [channel-type] --output-format [json|tree]",
		Short: "List channel types",
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
