package cli

import (
	"bytes"
	"context"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/stretchr/testify/require"
)

func TestCreateChannel(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	t.Cleanup(func() {
		deleteChannel(ch)
	})

	c := initClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, ch, resp.Channel.ID)
}

func TestCreateChannelAlreadyExists(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	t.Cleanup(func() {
		deleteChannel(ch)
	})
	time.Sleep(4 * time.Second)
	_ = app.Run([]string{"", "channel", "create", "-t", "messaging", "-i", ch, "-u", "userid"})
	require.Contains(t, app.ErrWriter.(*bytes.Buffer).String(), "channel exists already")
}

func TestGetChannel(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	t.Cleanup(func() {
		deleteChannel(ch)
	})

	err := app.Run([]string{"", "channel", "get", "-t", "messaging", "-i", ch})
	require.NoError(t, err)
}

func TestDeleteChannel(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	err := app.Run([]string{"", "channel", "delete", "-t", "messaging", "-i", ch, "--hard"})
	require.NoError(t, err)

	c := initClient()
	ctx := context.Background()
	_, err = c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.Error(t, err)
}

func TestUpdateChannel(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	t.Cleanup(func() {
		deleteChannel(ch)
	})

	err := app.Run([]string{"", "channel", "update", "-t", "messaging", "-i", ch, "-p", "{\"custom_property\":\"property-value\"}"})
	require.NoError(t, err)

	c := initClient()
	ctx := context.Background()
	resp, err := c.Channel("messaging", ch).Query(ctx, &stream.QueryRequest{Data: &stream.ChannelRequest{}})
	require.NoError(t, err)
	require.Equal(t, "property-value", resp.Channel.ExtraData["custom_property"])
}

func TestListChannel(t *testing.T) {
	app := initApp()

	err := app.Run([]string{"", "channel", "list", "-t", "messaging", "-l", "1"})
	require.NoError(t, err)
}
