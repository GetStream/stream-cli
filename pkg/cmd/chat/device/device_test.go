package device

import (
	"bytes"
	"testing"

	"github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
)

func TestDevice(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	userID := test.CreateUser()
	deviceID := test.RandomString(10)
	t.Cleanup(func() {
		test.DeleteUser(userID)
	})

	cmd.SetArgs([]string{"create-device", "-i", deviceID, "-p", "apn", "-u", userID})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	cmd.SetArgs([]string{"list-devices", "-u", userID})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), deviceID)

	cmd.SetArgs([]string{"delete-device", "-i", deviceID, "-u", userID})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully deleted device")
}
