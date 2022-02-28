package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateChannel(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	t.Cleanup(func() {
		app.Run([]string{"", "channel", "delete", "-t", "messaging", "-n", ch, "--hard"})
	})
}

func TestGetChannel(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	t.Cleanup(func() {
		app.Run([]string{"", "channel", "delete", "-t", "messaging", "-n", ch, "--hard"})
	})

	err := app.Run([]string{"", "channel", "get", "-t", "messaging", "-n", ch})
	require.NoError(t, err)
}

func TestDeleteChannel(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	err := app.Run([]string{"", "channel", "delete", "-t", "messaging", "-n", ch, "--hard"})
	require.NoError(t, err)
}

func TestUpdateChannel(t *testing.T) {
	app := initApp()
	ch := initChannel(t, app)
	t.Cleanup(func() {
		app.Run([]string{"", "channel", "delete", "-t", "messaging", "-n", ch, "--hard"})
	})

	err := app.Run([]string{"", "channel", "update", "-t", "messaging", "-n", ch, "-p", "{\"custom_property\":\"property-value\"}"})
	require.NoError(t, err)
}

func TestListChannel(t *testing.T) {
	app := initApp()

	err := app.Run([]string{"", "channel", "list", "-t", "messaging", "-l", "1"})
	require.NoError(t, err)
}
