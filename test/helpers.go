package test

import (
	"bytes"
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func prepareViperConfig() {
	viper.Set("default", "default_app")
	viper.Set("apps", []config.App{
		{
			Name:            "default_app",
			AccessKey:       os.Getenv("STREAM_KEY"),
			AccessSecretKey: os.Getenv("STREAM_SECRET"),
			URL:             config.DefaultEdgeURL,
		},
	})
}

func GetRootCmdWithSubCommands(c ...*cobra.Command) *cobra.Command {
	prepareViperConfig()

	rootCmd := &cobra.Command{}
	rootCmd.PersistentFlags().String("app", "", "app name")
	rootCmd.AddCommand(c...)
	rootCmd.SetIn(&bytes.Buffer{})
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})

	return rootCmd
}

func InitChannel(t *testing.T) string {
	name := RandomString(10)
	c := InitClient()
	c.CreateChannel(context.Background(), "messaging", name, "userid", nil)
	return name
}

func DeleteChannel(id string) {
	c := InitClient()
	_, _ = c.DeleteChannels(context.Background(), []string{"messaging:" + id}, true)
}

func InitClient() *stream.Client {
	c, _ := stream.NewClientFromEnvVars()
	return c
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) // A=65 and Z = 65+25
	}
	return string(bytes)
}
