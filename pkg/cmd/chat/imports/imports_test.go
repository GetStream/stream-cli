package imports

import (
	"bytes"
	"encoding/json"
	"testing"

	stream_chat "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
)

func TestValidateImport(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"validate-import", "./validator/testdata/valid-data.json"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
}

func TestUploadImport(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"upload-import", "./validator/testdata/valid-data.json"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
}

func TestGetImport(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"upload-import", "./validator/testdata/valid-data.json"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	var task stream_chat.ImportTask
	require.NoError(t, json.Unmarshal(cmd.OutOrStdout().(*bytes.Buffer).Bytes(), &task))

	cmd.SetArgs([]string{"get-import", task.ID})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
}

func TestListImports(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	cmd.SetArgs([]string{"list-imports", "--limit", "1"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)

	var tasks []stream_chat.ImportTask
	require.NoError(t, json.Unmarshal(cmd.OutOrStdout().(*bytes.Buffer).Bytes(), &tasks))

	require.Len(t, tasks, 1)
}
