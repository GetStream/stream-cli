package user

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GetStream/stream-cli/test"
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

func TestDeleteMultipleUsers(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	u1 := test.CreateUser()
	u2 := test.CreateUser()

	cmd.SetArgs([]string{"delete-users", "--hard-delete-users", "--hard-delete-messages", "--hard-delete-conversations", u1, u2})
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

func TestRevokeToken(t *testing.T) {
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(u)
	})
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"revoke-token", "-u", u})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully revoked token")
}

func TestBanUnbanUser(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	userToBan := test.CreateUser()
	bannerUser := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(userToBan)
		test.DeleteUser(bannerUser)
	})

	cmd.SetArgs([]string{"ban-user", "-t", userToBan, "-b", bannerUser, "--reason", "test"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully banned user")

	cmd.SetArgs([]string{"unban-user", "-t", userToBan})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully unbanned user")
}

func TestDeactivateReactivate(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	userToDeactivate := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(userToDeactivate)
	})

	cmd.SetArgs([]string{"deactivate-user", "-u", userToDeactivate})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully deactivated user")

	cmd.SetArgs([]string{"reactivate-user", "-u", userToDeactivate})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully reactivated user")
}

func TestMuteUnmute(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	userToMute := test.CreateUser()
	muterUser := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(userToMute)
		test.DeleteUser(muterUser)
	})

	cmd.SetArgs([]string{"mute-user", "-t", userToMute, "-b", muterUser, "--expiration", "5"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully muted user")

	cmd.SetArgs([]string{"unmute-user", "-t", userToMute, "-b", muterUser})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully unmuted user")
}

func TestFlag(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	userToFlag := test.CreateUser()
	flaggerUser := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(userToFlag)
		test.DeleteUser(flaggerUser)
	})

	cmd.SetArgs([]string{"flag-user", "-u", userToFlag, "-b", flaggerUser})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully flagged user")
}
