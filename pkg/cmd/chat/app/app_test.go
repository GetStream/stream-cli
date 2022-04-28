package app

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GetStream/stream-cli/test"
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

func TestRevokeAlTokens(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"revoke-all-tokens", "--before", "2000"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully revoked all tokens")
}
