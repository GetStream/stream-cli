package role

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listCmd(),
		createCmd(),
		deleteCmd(),
	}
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-roles --output-format [json|tree]",
		Short: "List all available roles",
		Example: heredoc.Doc(`
			# List all roles
			$ stream-cli chat list-roles
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.Permissions().ListRoles(cmd.Context())
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
		Use:   "create-role --name [role-name]",
		Short: "Create a custom role",
		Example: heredoc.Doc(`
			# Create a role named 'moderator-plus'
			$ stream-cli chat create-role --name moderator-plus
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")

			_, err = c.Permissions().CreateRole(cmd.Context(), name)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully created role [%s]\n", name)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("name", "n", "", "[required] Role name")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-role [name]",
		Short: "Delete a custom role",
		Example: heredoc.Doc(`
			# Delete the 'moderator-plus' role
			$ stream-cli chat delete-role moderator-plus
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			_, err = c.Permissions().DeleteRole(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			cmd.Printf("Successfully deleted role [%s]\n", args[0])
			return nil
		},
	}

	return cmd
}
