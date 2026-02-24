package permission

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		listCmd(),
		getCmd(),
	}
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-permissions --output-format [json|tree]",
		Short: "List all available permissions",
		Example: heredoc.Doc(`
			# List all permissions
			$ stream-cli chat list-permissions
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.Permissions().ListPermissions(cmd.Context())
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
		Use:   "get-permission [permission-id] --output-format [json|tree]",
		Short: "Get a permission by ID",
		Example: heredoc.Doc(`
			# Get a specific permission
			$ stream-cli chat get-permission send-message
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			resp, err := c.Permissions().GetPermission(cmd.Context(), args[0])
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
