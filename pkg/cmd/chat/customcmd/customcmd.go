package customcmd

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
		Use:   "list-commands --output-format [json|tree]",
		Short: "List all custom commands",
		Example: heredoc.Doc(`
			# List all custom commands
			$ stream-cli chat list-commands
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.ListCommands(cmd.Context())
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
		Use:   "get-command [name] --output-format [json|tree]",
		Short: "Get a custom command by name",
		Example: heredoc.Doc(`
			# Get the 'giphy' command
			$ stream-cli chat get-command giphy
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.GetCommand(cmd.Context(), args[0])
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
		Use:   "create-command --name [name] --description [description]",
		Short: "Create a custom command",
		Example: heredoc.Doc(`
			# Create a custom command named 'ticket'
			$ stream-cli chat create-command --name ticket --description "Create a support ticket"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("description")
			cmdArgs, _ := cmd.Flags().GetString("args")
			set, _ := cmd.Flags().GetString("set")

			newCmd := &stream.Command{
				Name:        name,
				Description: desc,
				Args:        cmdArgs,
				Set:         set,
			}

			resp, err := c.CreateCommand(cmd.Context(), newCmd)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully created command [%s]\n", resp.Command.Name)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Command name")
	fl.StringP("description", "d", "", "[required] Command description")
	fl.StringP("args", "a", "", "[optional] Command arguments help text")
	fl.StringP("set", "s", "", "[optional] Command set name for grouping")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("description")

	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-command --name [name] --description [description]",
		Short: "Update a custom command",
		Example: heredoc.Doc(`
			# Update the 'ticket' command description
			$ stream-cli chat update-command --name ticket --description "Create a new support ticket"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("description")
			cmdArgs, _ := cmd.Flags().GetString("args")
			set, _ := cmd.Flags().GetString("set")

			update := &stream.Command{
				Description: desc,
				Args:        cmdArgs,
				Set:         set,
			}

			_, err = c.UpdateCommand(cmd.Context(), name, update)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully updated command [%s]\n", name)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Command name")
	fl.StringP("description", "d", "", "[required] Command description")
	fl.StringP("args", "a", "", "[optional] Command arguments help text")
	fl.StringP("set", "s", "", "[optional] Command set name for grouping")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("description")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-command [name]",
		Short: "Delete a custom command",
		Example: heredoc.Doc(`
			# Delete the 'ticket' command
			$ stream-cli chat delete-command ticket
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			_, err = c.DeleteCommand(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			cmd.Printf("Successfully deleted command [%s]\n", args[0])
			return nil
		},
	}

	return cmd
}
