package app

import (
	"encoding/json"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{getCmd(), updateCmd(), revokeAllTokensCmd()}
}

func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-app --output-format [json|tree]",
		Short: "Get application settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			r, err := c.GetAppConfig(cmd.Context())
			if err != nil {
				return err
			}

			return utils.PrintObject(cmd, r.App)
		},
	}

	fl := cmd.Flags()
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")

	return cmd
}

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-app --properties [raw-json-update-properties]",
		Short: "Update application settings",
		Example: heredoc.Doc(`
			update-app --properties '{"multi_tenant_enabled": true, "permission_version": "v2"}'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			p, _ := cmd.Flags().GetString("properties")

			s := &stream.AppSettings{}
			err = json.Unmarshal([]byte(p), s)
			if err != nil {
				return err
			}

			_, err = c.UpdateAppSettings(cmd.Context(), s)
			if err != nil {
				return err
			}

			cmd.Println("Successfully updated app settings.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Raw json properties to update")
	cmd.MarkFlagRequired("properties")

	return cmd
}

func revokeAllTokensCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-all-tokens --before [epoch]",
		Short: "Revoke all tokens",
		Long: heredoc.Doc(`
			This command revokes ALL tokens for all users of an application. 
			This should be used with caution as it will expire every user’s token,
			regardless of whether the token has an iat claim.
		`),
		Example: heredoc.Doc(`
			# Revoke all tokens for the default app, from now
			$ stream-cli chat revoke-all-tokens

			# Revoke all tokens for the test app, before 2019-01-01
			$ stream-cli chat revoke-all-tokens --before 1546300800 --app test
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			before, _ := cmd.Flags().GetInt64("before")
			beforeDate := time.Unix(before, 0)

			_, err = c.RevokeTokens(cmd.Context(), &beforeDate)
			if err != nil {
				return err
			}

			cmd.Println("Successfully revoked all tokens.")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.IntP("before", "b", int(time.Now().Unix()), "[optional] The epoch timestamp before which tokens should be revoked. Defaults to now.")

	return cmd
}
