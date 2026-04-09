package utils

import (
	"context"
	"net/http"
	"os"
)

// UploadToS3 uploads a local file to the given S3 presigned URL via HTTP PUT.
func UploadToS3(ctx context.Context, filename, url string) error {
	data, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer data.Close()

	stat, err := data.Stat()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = stat.Size()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
