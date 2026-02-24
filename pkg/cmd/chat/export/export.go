package export

import (
	"encoding/json"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		exportChannelsCmd(),
		exportUserCmd(),
	}
}

func exportChannelsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-channels --channels [raw-json]",
		Short: "Export channel data to JSON",
		Long: heredoc.Doc(`
			Exports channel data to a JSON file. This is an async operation.
			A task ID will be returned which can be polled using the watch command.
		`),
		Example: heredoc.Doc(`
			# Export a specific channel
			$ stream-cli chat export-channels --channels '[{"type":"messaging","id":"general"}]'

			# Wait for the export to complete
			$ stream-cli chat watch <task-id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			channelsStr, _ := cmd.Flags().GetString("channels")

			var channels []*stream.ExportableChannel
			if err := json.Unmarshal([]byte(channelsStr), &channels); err != nil {
				return err
			}

			resp, err := c.ExportChannels(cmd.Context(), channels, nil)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully initiated channel export. Task id: %s\n", resp.TaskID)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("channels", "c", "", "[required] JSON array of channels to export")
	_ = cmd.MarkFlagRequired("channels")

	return cmd
}

func exportUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-user [user-id] --output-format [json|tree]",
		Short: "Export a user's data",
		Long: heredoc.Doc(`
			Exports a user's profile, reactions and messages. Raises an error
			if a user has more than 10k messages or reactions.
		`),
		Example: heredoc.Doc(`
			# Export data for user 'joe'
			$ stream-cli chat export-user joe
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.ExportUser(cmd.Context(), args[0])
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
