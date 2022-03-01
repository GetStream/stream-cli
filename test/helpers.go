package test

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/GetStream/stream-cli/pkg/config"
)

func InitConfig() *config.Config {
	return &config.Config{
		Default: "app",
		Apps: []config.App{
			{
				Name:            "app",
				AccessKey:       os.Getenv("STREAM_KEY"),
				AccessSecretKey: os.Getenv("STREAM_SECRET"),
			},
		},
		FilePath: "",
	}
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
