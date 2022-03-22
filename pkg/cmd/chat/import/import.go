package imports

import (
	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	cmd := &cobra.Command{
		Use:   "imports [subcommand]",
		Short: "Manage Chat imports",
		//Example:
	}

	cmd.AddCommand()
	return []*cobra.Command{cmd}
}

func validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import validate --filename [filename]",
		Short: "Validate an imports file",
		Example: heredoc.Doc(`
			# Validate imports file name 'imports.json'
			$ stream-cli chat imports validate --filename imports.json"

		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			filename, _ := cmd.Flags().GetString("filename")
			results, err := ValidateFile(cmd.Context(), c, filename)
			if err != nil {
				return err
			}

			PrintValidationResults(c.App.Writer, results)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("filename", "f", "", "[required] file name to be validated")
	cmd.MarkFlagRequired("filename")
	return cmd
	/*
		return &cli.Command{
				Name:  "validate",
				Usage: "validate an imports file",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "filename", Aliases: []string{"f"}, Required: true},
				},
				Action: func(c *cli.Context) error {
					state, err := chat.NewState(cfg)
					if err != nil {
						return err
					}
					filename := c.String("filename")
					if err := state.ValidateFile(c.Context, filename); err != nil {
						return err
					}

					return nil
				},
			}
	*/
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-channel --type [channel-type] --id [channel-id]",
		Short: "Return a channel",
		Example: heredoc.Doc(`
			# Returns 'redteam' channel of 'messaging' channel type as JSON
			$ stream-cli chat get-channel --type messaging --id redteam

			# Returns 'blueteam' channel of 'messaging' channel type as a browsable tree
			$ stream-cli chat get-channel --type messaging --id blueteam --output-format tree
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			chanType, _ := cmd.Flags().GetString("type")
			chanId, _ := cmd.Flags().GetString("id")

			r, err := c.Channel(chanType, chanId).Query(cmd.Context(), &stream.QueryRequest{})
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, r.Channel)
		},
	}

	fl := cmd.Flags()
	fl.StringP("type", "t", "", "[required] Channel type such as 'messaging' or 'livestream'")
	fl.StringP("id", "i", "", "[required] Channel id")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("id")

	return cmd
}
