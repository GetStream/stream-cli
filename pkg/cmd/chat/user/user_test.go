package user

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(u)
	})

	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"create-token", "-u", u, "-e", strconv.FormatInt(time.Now().Unix(), 10)})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Token for user")
}

func TestUpsertUser(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	userID := test.RandomString(10)
	t.Cleanup(func() {
		test.DeleteUser(userID)
	})

	cmd.SetArgs([]string{"upsert-user", "-p", "{\"id\": \"" + userID + "\"}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully upserted user")
}

func TestDeleteUser(t *testing.T) {
	u := test.CreateUser()
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"delete-user", "-u", u, "--hard-delete", "--mark-messages-deleted", "--delete-conversations"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully initiated user deletion")
}

func TestQueryUser(t *testing.T) {
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(u)
	})
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"query-users", "-f", "{\"id\": \"" + u + "\"}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), u)
}
