package channel

import (
	"bytes"
	"context"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	test "github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func getTestApp() *cli.App {
	return &cli.App{
		Commands:  []*cli.Command{NewChannelCmd(test.InitConfig())},
		Writer:    &bytes.Buffer{},
		ErrWriter: &bytes.Buffer{},
		ExitErrHandler: func(context *cli.Context, err error) {
			context.App.ErrWriter.Write([]byte(err.Error()))
		},
	}
}

func TestCreateChannel(t *testing.T) {
	ch := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	c := test.InitClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, ch, resp.Channel.ID)
}

func TestCreateChannelAlreadyExists(t *testing.T) {
	app := getTestApp()
	ch := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})
	time.Sleep(4 * time.Second)
	_ = app.Run([]string{"", "channel", "create", "-t", "messaging", "-i", ch, "-u", "userid"})
	require.Contains(t, app.ErrWriter.(*bytes.Buffer).String(), "channel exists already")
}

func TestGetChannel(t *testing.T) {
	app := cli.App{Commands: []*cli.Command{NewChannelCmd(test.InitConfig())}}
	ch := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	err := app.Run([]string{"", "channel", "get", "-t", "messaging", "-i", ch})
	require.NoError(t, err)
}

func TestDeleteChannel(t *testing.T) {
	app := getTestApp()
	ch := test.InitChannel(t)
	err := app.Run([]string{"", "channel", "delete", "-t", "messaging", "-i", ch, "--hard"})
	require.NoError(t, err)

	c := test.InitClient()
	ctx := context.Background()
	_, err = c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.Error(t, err)
}

func TestUpdateChannel(t *testing.T) {
	app := getTestApp()
	ch := test.InitChannel(t)
	t.Cleanup(func() {
		test.DeleteChannel(ch)
	})

	err := app.Run([]string{"", "channel", "update", "-t", "messaging", "-i", ch, "-p", "{\"custom_property\":\"property-value\"}"})
	require.NoError(t, err)

	c := test.InitClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, "property-value", resp.Channel.ExtraData["custom_property"])
}

func TestListChannel(t *testing.T) {
	app := getTestApp()

	err := app.Run([]string{"", "channel", "list", "-t", "messaging", "-l", "1"})
	require.NoError(t, err)
}
