package push

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GetStream/stream-cli/test"
)

func TestPushTest(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	msgID := test.CreateMessage(ch, u)
	t.Cleanup(func() {
		test.DeleteMessage(msgID)
		test.DeleteChannel(ch)
		test.DeleteUser(u)
	})

	cmd.SetArgs([]string{"test-push", "--message-id", msgID, "--user-id", u, "--skip-devices", "true"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), msgID)
}
