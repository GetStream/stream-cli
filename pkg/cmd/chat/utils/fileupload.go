package utils

import (
	"os"
	"path/filepath"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/spf13/cobra"
)

type uploadType string

const (
	uploadTypeFile  uploadType = "file"
	uploadTypeImage uploadType = "image"
)

func UploadFile(c *stream.Client, cmd *cobra.Command, chType, chID, userID, filePath string) (string, error) {
	return uploadFile(c, cmd, uploadTypeFile, chType, chID, userID, filePath)
}

func UploadImage(c *stream.Client, cmd *cobra.Command, chType, chID, userID, filePath string) (string, error) {
	return uploadFile(c, cmd, uploadTypeImage, chType, chID, userID, filePath)
}

func uploadFile(c *stream.Client, cmd *cobra.Command, uploadtype uploadType, chType, chID, userID, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	req := stream.SendFileRequest{
		User:     &stream.User{ID: userID},
		FileName: filepath.Base(file.Name()),
		Reader:   file,
	}

	var resp *stream.SendFileResponse

	if uploadtype == uploadTypeImage {
		resp, err = c.Channel(chType, chID).SendImage(cmd.Context(), req)
	} else {
		resp, err = c.Channel(chType, chID).SendFile(cmd.Context(), req)
	}

	if err != nil {
		return "", err
	}

	return resp.File, nil
}
