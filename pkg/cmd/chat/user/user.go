package user

import (
	"encoding/json"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		createTokenCmd(),
		upsertCmd(),
		deleteCmd(),
		queryCmd()}
}

func createTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token --user [user-id] --expiration [epoch]",
		Short: "Create a token",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, secret, err := config.GetConfig(cmd).GetCredentials(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user")
			exp, _ := cmd.Flags().GetInt("expiration")

			claims := jwt.MapClaims{
				"user_id": userID,
			}
			if exp > 0 {
				claims["exp"] = exp
			}

			token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
			if err != nil {
				return err
			}

			cmd.Printf("Token for user [%s]:\n%s\n", userID, token)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("user", "u", "", "[required] Id of the user to create token for")
	fl.IntP("expiration", "e", 0, "[optional] Expiration of the JWT in epoch seconds")
	cmd.MarkFlagRequired("user")

	return cmd
}

func upsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-user --properties [raw-json]",
		Short: "Upsert a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			props, _ := cmd.Flags().GetString("properties")

			user := &stream.User{}
			err = json.Unmarshal([]byte(props), user)
			if err != nil {
				return err
			}

			_, err = c.UpsertUser(cmd.Context(), user)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully upserted user [%s]\n", "")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("properties", "p", "", "[required] Raw JSON properties of the user")
	cmd.MarkFlagRequired("properties")

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-user --user [user-id] --hard-delete [true|false] --mark-messages-deleted [true|false] --delete-conversations [true|false]",
		Short: "Delete a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user")
			hardDelete, _ := cmd.Flags().GetBool("hard-delete")
			markMessagesDeleted, _ := cmd.Flags().GetBool("mark-messages-deleted")
			deleteConversations, _ := cmd.Flags().GetBool("delete-conversations")

			opts := []stream.DeleteUserOption{}
			if hardDelete {
				opts = append(opts, stream.DeleteUserWithHardDelete())
			}
			if markMessagesDeleted {
				opts = append(opts, stream.DeleteUserWithMarkMessagesDeleted())
			}
			if deleteConversations {
				opts = append(opts, stream.DeleteUserWithDeleteConversations())
			}

			_, err = c.DeleteUser(cmd.Context(), userID, opts...)
			if err != nil {
				return err
			}

			cmd.Println("Successfully deleted user")
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("user", "u", "", "[required] Id of the user to delete")
	fl.Bool("hard-delete", false, "[optional] Delete the user permanently")
	fl.Bool("mark-messages-deleted", false, "[optional] Mark messages of the user as deleted")
	fl.Bool("delete-conversations", false, "[optional] Delete all conversations of the user")
	cmd.MarkFlagRequired("user")

	return cmd
}

func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-users --filter [raw-json] --limit [int] --output-format [json|tree]",
		Short: "Query users",
		Example: heredoc.Doc(`
			query-users --filter '{"id": {"$eq": "user-1"}}' --limit 10 --output-format json
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			limit, _ := cmd.Flags().GetInt("limit")
			filter, _ := cmd.Flags().GetString("filter")

			var m map[string]interface{}
			err = json.Unmarshal([]byte(filter), &m)
			if err != nil {
				return err
			}

			q := &stream.QueryOption{
				Filter: m,
				Limit:  limit,
			}
			resp, err := c.QueryUsers(cmd.Context(), q)
			if err != nil {
				return err
			}

			utils.PrintObject(cmd, resp)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("filter", "f", "{}", "[required] Filter for users")
	fl.IntP("limit", "l", 10, "[optional] The number of users returned")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	cmd.MarkFlagRequired("filter")

	return cmd
}
