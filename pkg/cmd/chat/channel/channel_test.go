package channel

import (
	"bytes"
	"context"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
)

func TestCreateChannel(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.RandomString(10)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"create-channel", "-t", "messaging", "-i", ch, "-u", "user"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	c := test.InitClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, ch, resp.Channel.ID)
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

	cmd.SetArgs([]string{"update-channel-partial", "-t", "messaging", "-i", ch, "-s", "color=blue,age=27"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	c := test.InitClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, "blue", resp.Channel.ExtraData["color"])
	require.Equal(t, "27", resp.Channel.ExtraData["age"])
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
