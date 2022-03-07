package channeltype

import (
	"bytes"
	"encoding/json"
	"fmt"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
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
		Use:   "get-channel-type --channel-type [channel-type] --output-format [json]",
		Short: "Get channel type",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			ct, _ := cmd.Flags().GetString("channel-type")
			format, _ := cmd.Flags().GetString("output-format")

			resp, err := c.GetChannelType(cmd.Context(), ct)
			if err != nil {
				return err
			}

			if format == "json" {
				unindented, err := json.Marshal(resp)
				if err != nil {
					return err
				}

				var indented bytes.Buffer
				err = json.Indent(&indented, unindented, "", "  ")
				if err != nil {
					return err
				}

				cmd.Println(indented.String())
			} else {
				return fmt.Errorf("unknown output format: %s", format)
			}

			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channel-type", "t", "", "Channel type")
	fl.StringP("output-format", "o", "json", "Output format. Can be json or [see-in-next-pull-request]")
	cmd.MarkFlagRequired("channel-type")

	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-channel-type --properties [raw-json-properties]",
		Short: "Create channel type",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
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
	fl.StringP("properties", "p", "", "Raw JSON properties")
	cmd.MarkFlagRequired("properties")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-channel-type --channel-type [channel-type]",
		Short: "Delete channel type",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			ct, _ := cmd.Flags().GetString("channel-type")

			_, err = c.GetChannelType(cmd.Context(), ct)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted channel type.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channel-type", "t", "", "Channel type")
	cmd.MarkFlagRequired("channel-type")

	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-channel-type --channel-type [channel-type] --properties [raw-json-properties]",
		Short: "Update channel type",
		Example: heredoc.Doc(`
			update-channel-type --channel-type my-channel-type --properties '{"quotes": true}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			ct, _ := cmd.Flags().GetString("channel-type")
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
	fl.StringP("channel-type", "t", "", "Channel type")
	fl.StringP("properties", "p", "", "Raw JSON properties")
	cmd.MarkFlagRequired("channel-type")
	cmd.MarkFlagRequired("properties")

	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-channel-types --channel-type [channel-type] --output-format [json]",
		Short: "List channel types",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetStreamClient(cmd)
			if err != nil {
				return err
			}

			format, _ := cmd.Flags().GetString("output-format")

			resp, err := c.ListChannelTypes(cmd.Context())
			if err != nil {
				return err
			}

			if format == "json" {
				unindented, err := json.Marshal(resp)
				if err != nil {
					return err
				}

				var indented bytes.Buffer
				err = json.Indent(&indented, unindented, "", "  ")
				if err != nil {
					return err
				}

				cmd.Println(indented.String())
			} else {
				return fmt.Errorf("unknown output format: %s", format)
			}

			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "Output format. Can be json or [see-in-next-pull-request]")

	return cmd
}
