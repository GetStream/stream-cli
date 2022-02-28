package cli

import (
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func initApp() *cli.App {
	c := &Config{
		Default: "app",
		Apps: []App{
			{
				Name:            "app",
				AccessKey:       os.Getenv("STREAM_KEY"),
				AccessSecretKey: os.Getenv("STREAM_SECRET"),
			},
		},
		filePath: "",
	}

	return NewRootCmd(c)
}

func initChannel(t *testing.T, app *cli.App) string {
	name := randomString(10)
	err := app.Run([]string{"", "channel", "create", "-t", "messaging", "-n", name, "-u", "userid"})
	require.NoError(t, err)
	return name
}

func randomString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) // A=65 and Z = 65+25
	}
	return string(bytes)
}
