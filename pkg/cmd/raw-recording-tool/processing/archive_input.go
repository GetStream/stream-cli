package processing

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/GetStream/getstream-go/v3"
)

// extractToTempDir extracts archive to temp directory or returns the directory path
// Returns: (workingDir, cleanupFunc, error)
func ExtractToTempDir(inputPath string, logger *getstream.DefaultLogger) (string, func(), error) {
	// If it's already a directory, just return it
	if stat, err := os.Stat(inputPath); err == nil && stat.IsDir() {
		logger.Debug("Input is already a directory: %s", inputPath)
		return inputPath, func() {}, nil
	}

	// If it's a tar.gz file, extract it to temp directory
	if strings.HasSuffix(strings.ToLower(inputPath), ".tar.gz") {
		logger.Info("Extracting tar.gz archive to temporary directory...")

		tempDir, err := os.MkdirTemp("", "raw-tools-*")
		if err != nil {
			return "", nil, fmt.Errorf("failed to create temp directory: %w", err)
		}

		cleanup := func() {
			os.RemoveAll(tempDir)
		}

		err = extractTarGzToDir(inputPath, tempDir, logger)
		if err != nil {
			cleanup()
			return "", nil, fmt.Errorf("failed to extract tar.gz: %w", err)
		}

		logger.Debug("Extracted archive to: %s", tempDir)
		return tempDir, cleanup, nil
	}

	return "", nil, fmt.Errorf("unsupported input format: %s (only tar.gz files and directories supported)", inputPath)
}

// extractTarGzToDir extracts a tar.gz file to the specified directory
func extractTarGzToDir(tarGzPath, destDir string, logger *getstream.DefaultLogger) error {
	file, err := os.Open(tarGzPath)
	if err != nil {
		return fmt.Errorf("failed to open tar.gz file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Skip directories
		if header.FileInfo().IsDir() {
			continue
		}

		// Create destination file
		destPath := filepath.Join(destDir, header.Name)

		// Create directory structure if needed
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory structure: %w", err)
		}

		// Extract file
		outFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", destPath, err)
		}

		_, err = io.Copy(outFile, tarReader)
		outFile.Close()
		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", destPath, err)
		}
	}

	return nil
}
