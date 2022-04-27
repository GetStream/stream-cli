package channeltype

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/GetStream/stream-cli/test"
)

func deleteChannelType(name string) {
	c := test.InitClient()

	for i := 0; i < 5; i++ {
		_, err := c.DeleteChannelType(context.Background(), name)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
}

func TestCreateChannelType(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	name := test.RandomString(10)
	t.Cleanup(func() {
		deleteChannelType(name)
	})

	cmd.SetArgs([]string{"create-channel-type", "-p", "{\"name\":\"" + name + "\"}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
}

func TestUpdateChannelType(t *testing.T) {
	t.Skip("Fix this")
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	name := test.RandomString(10)
	t.Cleanup(func() {
		deleteChannelType(name)
	})

	cmd.SetArgs([]string{"create-channel-type", "-p", "{\"name\":\"" + name + "\"}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	cmd.SetArgs([]string{"update-channel-type", "-t", name, "-p", "{\"quotes\": true}"})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
}

func TestDeleteChannelType(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	name := test.RandomString(10)
	t.Cleanup(func() {
		deleteChannelType(name)
	})

	cmd.SetArgs([]string{"create-channel-type", "-p", "{\"name\":\"" + name + "\"}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	cmd.SetArgs([]string{"delete-channel-type", name})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
}

func TestListChannelTypes(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)

	cmd.SetArgs([]string{"list-channel-types"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "messaging")
}
