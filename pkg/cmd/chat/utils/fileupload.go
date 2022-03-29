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

func UploadFile(c *stream.Client, cmd *cobra.Command, chType, chId, userID, filePath, contentType string) (string, error) {
	return uploadFileInternal(c, cmd, uploadTypeFile, chType, chId, userID, filePath, contentType)
}

func UploadImage(c *stream.Client, cmd *cobra.Command, chType, chId, userID, filePath, contentType string) (string, error) {
	return uploadFileInternal(c, cmd, uploadTypeImage, chType, chId, userID, filePath, contentType)
}

func uploadFileInternal(c *stream.Client, cmd *cobra.Command, uploadtype uploadType, chType, chId, userID, filePath, contentType string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	req := stream.SendFileRequest{
		User:        &stream.User{ID: userID},
		FileName:    filepath.Base(file.Name()),
		ContentType: contentType,
		Reader:      file,
	}

	var resp *stream.SendFileResponse

	if uploadtype == uploadTypeImage {
		resp, err = c.Channel(chType, chId).SendImage(cmd.Context(), req)
	} else {
		resp, err = c.Channel(chType, chId).SendFile(cmd.Context(), req)
	}

	if err != nil {
		return "", err
	}

	return resp.File, nil
}

func DeleteFile(c *stream.Client, cmd *cobra.Command, chType, chId, fileUrl string) error {
	_, err := c.Channel(chType, chId).DeleteFile(cmd.Context(), fileUrl)
	return err
}

func DeleteImage(c *stream.Client, cmd *cobra.Command, chType, chId, fileUrl string) error {
	_, err := c.Channel(chType, chId).DeleteImage(cmd.Context(), fileUrl)
	return err
}
