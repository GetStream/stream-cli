package blocklist

import (
	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listCmd(),
		getCmd(),
		createCmd(),
		updateCmd(),
		deleteCmd(),
	}
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-blocklists --output-format [json|tree]",
		Short: "List all block lists",
		Example: heredoc.Doc(`
			# List all block lists
			$ stream-cli chat list-blocklists
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.ListBlocklists(cmd.Context())
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

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-blocklist [name] --output-format [json|tree]",
		Short: "Get a block list by name",
		Example: heredoc.Doc(`
			# Get the 'profanity' block list
			$ stream-cli chat get-blocklist profanity
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.GetBlocklist(cmd.Context(), args[0])
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
		Use:   "create-blocklist --name [name] --words [comma-separated-words]",
		Short: "Create a block list",
		Example: heredoc.Doc(`
			# Create a block list named 'profanity' with some words
			$ stream-cli chat create-blocklist --name profanity --words "badword1,badword2,badword3"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			words, err := utils.GetStringSliceParam(cmd.Flags(), "words")
			if err != nil {
				return err
			}

			bl := &stream.BlocklistCreateRequest{
				BlocklistBase: stream.BlocklistBase{
					Name:  name,
					Words: words,
				},
			}

			_, err = c.CreateBlocklist(cmd.Context(), bl)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully created block list [%s]\n", name)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Block list name")
	fl.StringP("words", "w", "", "[required] Comma-separated list of words to block")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("words")

	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-blocklist --name [name] --words [comma-separated-words]",
		Short: "Update a block list",
		Example: heredoc.Doc(`
			# Update the 'profanity' block list
			$ stream-cli chat update-blocklist --name profanity --words "newword1,newword2"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			words, err := utils.GetStringSliceParam(cmd.Flags(), "words")
			if err != nil {
				return err
			}

			_, err = c.UpdateBlocklist(cmd.Context(), name, words)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully updated block list [%s]\n", name)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Block list name")
	fl.StringP("words", "w", "", "[required] Comma-separated list of words to block")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("words")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-blocklist [name]",
		Short: "Delete a block list",
		Example: heredoc.Doc(`
			# Delete the 'profanity' block list
			$ stream-cli chat delete-blocklist profanity
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			_, err = c.DeleteBlocklist(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			cmd.Printf("Successfully deleted block list [%s]\n", args[0])
			return nil
		},
	}

	return cmd
}
