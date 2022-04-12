package message

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	t.Cleanup(func() {
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"send-message", "-t", "messaging", "-i", ch, "-u", u, "--attachment", "https://via.placeholder.com/1", "--text", "hello"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Message successfully sent")
}

func TestSendMessageWithFileAttachment(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()

	tmpfile, err := ioutil.TempFile("", "*.txt")
	require.NoError(t, err)

	err = os.WriteFile(tmpfile.Name(), []byte("hello\nworld\n"), 0644)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name())
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"send-message", "-t", "messaging", "-i", ch, "-u", u, "--attachment", tmpfile.Name(), "--text", "hello"})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Message successfully sent")
}

func TestGetMessage(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	m := test.CreateMessage(ch, u)
	t.Cleanup(func() {
		test.DeleteMessage(m)
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"get-message", m})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), m)
}

func TestGetMultipleMessage(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	m1 := test.CreateMessage(ch, u)
	m2 := test.CreateMessage(ch, u)
	t.Cleanup(func() {
		test.DeleteMessage(m1)
		test.DeleteMessage(m2)
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"get-messages", "-t", "messaging", "-i", ch, m1, m2})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), m1)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), m2)
}

func TestDeleteMessage(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	m := test.CreateMessage(ch, u)
	t.Cleanup(func() {
		test.DeleteMessage(m)
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"delete-message", m})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Message successfully deleted")
}

func TestPartialUpdateMessage(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	m := test.CreateMessage(ch, u)
	t.Cleanup(func() {
		test.DeleteMessage(m)
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"update-message-partial", "-m", m, "--user", u, "--set", "age=15"})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully updated message")
}

func TestFlagMessage(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()
	m := test.CreateMessage(ch, u)
	t.Cleanup(func() {
		test.DeleteMessage(m)
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"flag-message", "-m", m, "-u", u})
	_, err := cmd.ExecuteC()
	require.NoError(t, err)
	require.Contains(t, cmd.OutOrStdout().(*bytes.Buffer).String(), "Successfully flagged message")
}
