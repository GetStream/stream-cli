package utils

import (
	"os"
	"path/filepath"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/spf13/cobra"
)

func UploadFile(c *stream.Client, cmd *cobra.Command, chType, chId, userID, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	req := stream.SendFileRequest{
		User:     &stream.User{ID: userID},
		FileName: filepath.Base(file.Name()),
		Reader:   file,
	}
	resp, err := c.Channel(chType, chId).SendFile(cmd.Context(), req)
	if err != nil {
		return "", err
	}

	return resp.File, nil
}
