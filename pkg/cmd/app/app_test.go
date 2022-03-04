package app

import (
	"bytes"
	"testing"

	"github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
)

func TestGetAppJsonFormat(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"get-app", "--output-format", "json"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "organization")
}

func TestGetAppUnknownFormat(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"get-app", "--output-format", "unknown"})
	_, err := cmd.ExecuteC()
	require.Error(t, err)
	require.NotContains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "organization")
}

func TestUpdateApp(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"update-app", "--properties", "{\"multi_tenant_enabled\":true}"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
}
