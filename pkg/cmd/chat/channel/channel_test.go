package channel

import (
	"bytes"
	"context"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/stretchr/testify/require"

	"github.com/GetStream/stream-cli/test"
)

func TestCreateChannel(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.RandomString(10)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"create-channel", "-t", "messaging", "-i", ch, "-u", "user", "-p", "{\"custom_property\":\"property-value\"}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	c := test.InitClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, ch, resp.Channel.ID)
	require.Equal(t, "property-value", resp.Channel.ExtraData["custom_property"])
}

func TestCreateChannelAlreadyExists(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	time.Sleep(4 * time.Second)
	cmd.SetArgs([]string{"create-channel", "-t", "messaging", "-i", ch, "-u", "user"})
	_, err := cmd.ExecuteC()
	require.Error(t, err)
	require.Contains(t, cmd.ErrOrStderr().(*bytes.Buffer).String(), "already exists")
}

func TestGetChannel(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"get-channel", "-t", "messaging", "-i", ch})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), ch)
}

func TestDeleteChannel(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	cmd.SetArgs([]string{"delete-channel", "-t", "messaging", "-i", ch, "--hard"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	c := test.InitClient()
	ctx := context.Background()
	_, err = c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.Error(t, err)
}

func TestUpdateChannel(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"update-channel", "-t", "messaging", "-i", ch, "-p", "{\"custom_property\":\"property-value\"}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	c := test.InitClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, "property-value", resp.Channel.ExtraData["custom_property"])
}

func TestUpdateChannelPartial(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"update-channel-partial", "-t", "messaging", "-i", ch, "-s", "{\"color\":\"blue\",\"age\":27}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	c := test.InitClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, "blue", resp.Channel.ExtraData["color"])
	require.Equal(t, 27., resp.Channel.ExtraData["age"])
}

func TestListChannel(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	chName := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(chName)
	})

	cmd.SetArgs([]string{"list-channels", "-t", "messaging", "-l", "1"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), chName)
}

func TestAddMembersRemoveMembers(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"add-members", "-t", "messaging", "-i", ch, u})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully added user")

	cmd.SetArgs([]string{"remove-members", "-t", "messaging", "-i", ch, u})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully removed user")
}

func TestPromoteAndDemoteModerator(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"add-members", "-t", "messaging", "-i", ch, u})
	_, _ = cmd.ExecuteC()

	cmd.SetArgs([]string{"promote-moderators", "-t", "messaging", "-i", ch, u})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully promoted user")

	cmd.SetArgs([]string{"demote-moderators", "-t", "messaging", "-i", ch, u})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully demoted user")
}

func TestAssignRole(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"add-members", "-t", "messaging", "-i", ch, u})
	_, _ = cmd.ExecuteC()

	cmd.SetArgs([]string{"assign-role", "-t", "messaging", "-i", ch, "-u", u, "-r", "channel_moderator"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully assigned role")
}

func TestHideAndShowChannel(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteChannel(ch)
		test.DeleteUser(u)
	})

	cmd.SetArgs([]string{"add-members", "-t", "messaging", "-i", ch, u})
	_, _ = cmd.ExecuteC()

	cmd.SetArgs([]string{"hide-channel", "-t", "messaging", "-i", ch, "-u", u})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully hid channel")

	cmd.SetArgs([]string{"show-channel", "-t", "messaging", "-i", ch, "-u", u})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully shown channel")
}

func TestListMembers(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteChannel(ch)
		test.DeleteUser(u)
	})

	// Add member so they show up in the member list
	cmd.SetArgs([]string{"add-members", "-t", "messaging", "-i", ch, u})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	// Now list members
	cmd.SetArgs([]string{"list-members", "-t", "messaging", "-i", ch})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)

	// Assert the user ID is in the output
	out := cmd.OutOrStdout().(*bytes.Buffer).String()
	require.Contains(t, out, u)
}
