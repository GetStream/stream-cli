package cli

import (
	"bytes"
	"context"
	"flag"
	"math/rand"
	"os"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func initConfig() *Config {
	return &Config{
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
}

func initApp() *cli.App {
	app := NewRootCmd(initConfig())

	// This way, it's easy to read the output of the app from tests.
	app.Writer = &bytes.Buffer{}
	app.ErrWriter = &bytes.Buffer{}
	app.Reader = &bytes.Buffer{}

	// We need to overwrite this because we can't just simply
	// exit the process during unit testing.
	app.ExitErrHandler = func(c *cli.Context, err error) {
		c.App.ErrWriter.Write([]byte(err.Error()))
	}

	return app
}

func initChannel(t *testing.T, app *cli.App) string {
	name := randomString(10)
	err := app.Run([]string{"", "channel", "create", "-t", "messaging", "-n", name, "-u", "userid"})
	require.NoError(t, err)
	return name
}

func deleteChannel(id string) {
	c := initClient()
	_, _ = c.DeleteChannels(context.Background(), []string{"messaging:" + id}, true)
}

func initClient() *stream.Client {
	c, _ := initConfig().GetStreamClient(cli.NewContext(nil, &flag.FlagSet{}, nil))
	return c
}

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) // A=65 and Z = 65+25
	}
	return string(bytes)
}
