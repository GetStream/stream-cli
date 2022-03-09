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
		queryCmd()}
}

func createTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-token --user [user-id] --expiration [epoch]",
		Short: "Create a token",
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
	fl.IntP("expiration", "e", 0, "[optional] Expiration (exp) of the JWT in epoch seconds")
	fl.IntP("issued-at", "i", 0, "[optional] Issued at (iat) of the JWT in epoch seconds")
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
	fl.Bool("hard-delete", false, "[optional] Delete the user permanently")
	fl.Bool("mark-messages-deleted", false, "[optional] Mark messages of the user as deleted")
	fl.Bool("delete-conversations", false, "[optional] Delete all conversations of the user")
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
