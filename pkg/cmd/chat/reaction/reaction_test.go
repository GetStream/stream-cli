package reaction

import (
	"bytes"
	"testing"

	"github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
)

func TestReactions(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	msgID := test.CreateMessage(ch, u)
	t.Cleanup(func() {
		test.DeleteMessage(msgID)
		test.DeleteChannel(ch)
		test.DeleteUser(u)
	})

	cmd.SetArgs([]string{"send-reaction", "-m", msgID, "-u", u, "-r", "like"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully sent reaction")

	cmd.SetArgs([]string{"get-reactions", msgID})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "like")

	cmd.SetArgs([]string{"delete-reaction", "-m", msgID, "-r", "like", "-u", u})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully deleted reaction")
}
