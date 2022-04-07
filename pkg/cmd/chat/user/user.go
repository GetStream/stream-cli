package user

import (
	"encoding/json"
	"errors"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/GetStream/stream-cli/pkg/utils"
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmds() []*cobra.Command {
	return []*cobra.Command{
		createTokenCmd(),
		upsertCmd(),
		deleteCmd(),
		queryCmd(),
		revokeCmd()}
}

func createTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token --user [user-id] --expiration [epoch] --issued-at [epoch]",
		Short: "Create a token",
		Long: heredoc.Doc(`
			Stream uses JWT (JSON Web Tokens) to authenticate chat users, enabling them to login.
			Knowing whether a user is authorized to perform certain actions is
			managed separately via a role based permissions system.

			With this command you can generate token for a specific user that can be
			used on the frontend.
		`),
		Example: heredoc.Doc(`
			# Create a JWT token for a user with id '123'. This token has no expiration.
			$ stream-cli chat create-token --user 123

			# Create a JWT for user 'joe' with 'exp' and 'iat' claim
			$ stream-cli chat create-token --user joe --expiration 1577880000 --issued-at 1577880000
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user")
			exp, _ := cmd.Flags().GetInt("expiration")
			iat, _ := cmd.Flags().GetInt("issued-at")

			expDate := time.Time{}
			iatDate := time.Time{}
			if exp > 0 {
				expDate = time.Unix(int64(exp), 0)
			}
			if iat > 0 {
				iatDate = time.Unix(int64(iat), 0)
			}

			token, err := c.CreateToken(userID, expDate, iatDate)
			if err != nil {
				return err
			}

			cmd.Printf("Token for user [%s]:\n%s\n", userID, token)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("user", "u", "", "[required] Id of the user to create token for")
	fl.IntP("expiration", "e", 0, "[optional] Expiration (exp) of the JWT in epoch timestamp")
	fl.IntP("issued-at", "i", 0, "[optional] Issued at (iat) of the JWT in epoch timestamp")
	cmd.MarkFlagRequired("user")

	return cmd
}

func upsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upsert-user --properties [raw-json]",
		Short: "Upsert a user",
		Long: heredoc.Doc(`
			This command inserts a new or updates an existing user.
			Stream Users require only an id to be created.
			Any user present in the payload will have its data replaced with the new version.
		`),
		Example: heredoc.Doc(`
			# Create a new user with id 'my-user-1'
			$ stream-cli chat upsert-user --properties "{\"id\":\"my-user-1\"}"

			Check the Go SDK's 'User' struct for the properties that you can use here.
		`),
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
		Long: heredoc.Doc(`
			This command deletes a user. If not flags are provided, user and messages will be soft deleted.
			
			There are 3 additional options that you can provide:

			--hard-delete: If set to true, hard deletes everything related to this user, channels, messages and everything related to it.
			--mark-messages-deleted: If set to true, hard deletes all messages related to this user.
			--delete-conversations: If set to true, hard deletes all conversations related to this user.

			User deletion is an async operation in the backend.
			Once it succeeded, you'll need to use the 'watch' command to see the async task's result.
		`),
		Example: heredoc.Doc(`
			# Soft delete a user with id 'my-user-1'
			$ stream-cli chat delete-user --user my-user-1

			# Hard delete a user with id 'my-user-2'
			$ stream-cli chat delete-user --user my-user-2 --hard-delete
			> Successfully initiated user deletion. Task id: 8d011daa-cbcd-4cba-ad16-701de599873a

			# Watch the async task's result
			$ stream-cli chat watch 8d011daa-cbcd-4cba-ad16-701de599873a
			> Async operation completed successfully.
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user")
			hardDelete, _ := cmd.Flags().GetBool("hard-delete")
			markMessagesDeleted, _ := cmd.Flags().GetBool("mark-messages-deleted")
			deleteConversations, _ := cmd.Flags().GetBool("delete-conversations")

			resp, err := c.DeleteUsers(cmd.Context(), []string{userID}, stream.DeleteUserOptions{
				User:          getDeleteType(hardDelete),
				Messages:      getDeleteType(markMessagesDeleted),
				Conversations: getDeleteType(deleteConversations)})
			if err != nil {
				return err
			}

			if resp.TaskID != "" {
				cmd.Printf("Successfully initiated user deletion. Task id: %s\n", resp.TaskID)
				return nil
			} else {
				return errors.New("user deletion failed")
			}
		},
	}

	fl := cmd.Flags()
	fl.StringP("user", "u", "", "[required] Id of the user to delete")
	fl.Bool("hard-delete", false, "[optional] Hard delete everything related to this user")
	fl.Bool("mark-messages-deleted", false, "[optional] Hard delete all messages related to the user")
	fl.Bool("delete-conversations", false, "[optional] Hard delete all conversations related to the user")
	cmd.MarkFlagRequired("user")

	return cmd
}

func getDeleteType(hardDeleteEnabled bool) stream.DeleteType {
	if hardDeleteEnabled {
		return stream.HardDelete
	}
	return stream.SoftDelete
}

func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-users --filter [raw-json] --limit [int] --output-format [json|tree]",
		Short: "Query users",
		Long: heredoc.Doc(`
			This command allows you to search for users. The 'filter' flag is a raw JSON string,
			and you can check the valid combinations in the official documentation.

			https://getstream.io/chat/docs/node/query_users/?language=javascript
		`),
		Example: heredoc.Doc(`
			# Query for 'user-1'. The results are shown as json.
			$ stream-cli chat query-users --filter '{"id": {"$eq": "user-1"}}'

			# Query for 'user-1' and 'user-2'. The results are shown as a browsable tree.
			$ stream-cli chat query-users --filter '{"id": {"$in": ["user-1", "user-2"]}}' --output-format tree
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

			return utils.PrintObject(cmd, resp)
		},
	}

	fl := cmd.Flags()
	fl.StringP("filter", "f", "{}", "[required] Filter for users")
	fl.IntP("limit", "l", 10, "[optional] The number of users returned")
	fl.StringP("output-format", "o", "json", "[optional] Output format. Can be json or tree")
	cmd.MarkFlagRequired("filter")

	return cmd
}

func revokeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-token --user [user-id] --before [epoch]",
		Short: "Revoke a token",
		Long: heredoc.Doc(`
			Revokes a token for a single user. All requests will be rejected that
			were issued before the given epoch timestamp.
		`),
		Example: heredoc.Doc(`
			# Revoke token for user 'joe' before today's date (default date)
			$ stream-cli revoke-token --user joe

			# Revoke token for user 'mike' before 2019-01-01
			$ stream-cli revoke-token --user mike --before 1546300800
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.GetConfig(cmd).GetClient(cmd)
			if err != nil {
				return err
			}

			userID, _ := cmd.Flags().GetString("user")
			before, _ := cmd.Flags().GetInt64("before")
			if before == 0 {
				before = time.Now().Unix()
			}
			beforeDate := time.Unix(before, 0)

			_, err = c.RevokeUserToken(cmd.Context(), userID, &beforeDate)
			if err != nil {
				return err
			}

			cmd.Printf("Successfully revoked token for user [%s]\n", userID)
			return nil
		},
	}

	fl := cmd.Flags()
	fl.StringP("user", "u", "", "[required] Id of the user to revoke token for")
	fl.Int64P("before", "b", 0, "[optional] The epoch timestamp before which tokens should be revoked. Defaults to now.")
	cmd.MarkFlagRequired("user")

	return cmd
}
