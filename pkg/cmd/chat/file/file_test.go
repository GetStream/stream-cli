package file

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/GetStream/stream-cli/test"
	"github.com/stretchr/testify/require"
)

func TestFileUploadAndDelete(t *testing.T) {
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

	cmd.SetArgs([]string{"upload-file", "-t", "messaging", "-i", ch, "-u", u, "--file", tmpfile.Name()})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
	stdOut := cmd.OutOrStdout().(*bytes.Buffer).String()
	require.Contains(t, stdOut, "Successfully uploaded file")

	// The stdout looks like this:
	// Successfully uploaded file: http://example.org/snippet.txt\n
	url := strings.Split(stdOut, ": ")[1]
	url = strings.TrimSuffix(url, "\n")
	cmd.SetArgs([]string{"delete-file", "-t", "messaging", "-i", ch, "-u", u, "--file-url", url})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
}

func TestImageUploadAndDelete(t *testing.T) {
	cmd := test.GetRootCmdWithSubCommands(NewCmds()...)
	ch := test.InitChannel(t)
	u := test.CreateUser()

	tmpfile, err := ioutil.TempFile("", "*.png")
	require.NoError(t, err)

	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	err = png.Encode(tmpfile, m)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name())
		test.DeleteUser(u)
		test.DeleteChannel(ch)
	})

	cmd.SetArgs([]string{"upload-image", "-t", "messaging", "-i", ch, "-u", u, "--file", tmpfile.Name()})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)

	stdOut := cmd.OutOrStdout().(*bytes.Buffer).String()
	require.Contains(t, stdOut, "Successfully uploaded image")

	// The stdout looks like this:
	// Successfully uploaded image: http://example.org/image.png\n
	url := strings.Split(stdOut, ": ")[1]
	url = strings.TrimSuffix(url, "\n")
	cmd.SetArgs([]string{"delete-image", "-t", "messaging", "-i", ch, "--image-url", url})
	_, err = cmd.ExecuteC()
	require.NoError(t, err)
}
