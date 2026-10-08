package rawrecording

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Downloader handles downloading files from S3 with caching
type S3Downloader struct {
	cacheDir string
	verbose  bool
}

// CacheMetadata stores information about a cached file
type CacheMetadata struct {
	ETag         string `json:"etag"`
	OriginalURL  string `json:"original_url"`
	LastModified string `json:"last_modified,omitempty"`
}

// NewS3Downloader creates a new S3Downloader
func NewS3Downloader(cacheDir string, verbose bool) *S3Downloader {
	return &S3Downloader{
		cacheDir: cacheDir,
		verbose:  verbose,
	}
}

// Download downloads a file from S3 or presigned URL, using cache if available
// Returns the local file path to the downloaded file
func (d *S3Downloader) Download(ctx context.Context, inputURL string) (string, error) {
	// Ensure cache directory exists
	if err := os.MkdirAll(d.cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Generate cache key from URL
	cacheKey := d.generateCacheKey(inputURL)
	cachedFilePath := filepath.Join(d.cacheDir, cacheKey+".tar.gz")
	metadataPath := filepath.Join(d.cacheDir, cacheKey+".meta.json")

	// Check if file is already cached
	if d.isCacheValid(ctx, inputURL, cachedFilePath, metadataPath) {
		if d.verbose {
			fmt.Printf("Using cached file: %s\n", cachedFilePath)
		}
		return cachedFilePath, nil
	}

	// Download the file
	if d.verbose {
		fmt.Printf("Downloading from: %s\n", d.sanitizeURLForLog(inputURL))
	}

	var etag string
	var err error

	if isS3URL(inputURL) {
		etag, err = d.downloadFromS3(ctx, inputURL, cachedFilePath)
	} else {
		etag, err = d.downloadFromPresignedURL(ctx, inputURL, cachedFilePath)
	}

	if err != nil {
		return "", err
	}

	// Save cache metadata
	metadata := CacheMetadata{
		ETag:        etag,
		OriginalURL: d.hashURL(inputURL), // Store hash instead of URL for privacy
	}
	if err := d.saveCacheMetadata(metadataPath, &metadata); err != nil {
		// Log but don't fail - download succeeded
		if d.verbose {
			fmt.Printf("Warning: failed to save cache metadata: %v\n", err)
		}
	}

	if d.verbose {
		fmt.Printf("Downloaded to: %s\n", cachedFilePath)
	}

	return cachedFilePath, nil
}

// generateCacheKey creates a unique cache key from the URL
func (d *S3Downloader) generateCacheKey(inputURL string) string {
	return d.hashURL(inputURL)
}

// hashURL creates a SHA256 hash of the URL
func (d *S3Downloader) hashURL(inputURL string) string {
	// For presigned URLs, we only hash the base path (without query params)
	// This allows the same file to be cached even if the signature changes
	baseURL := inputURL
	if u, err := url.Parse(inputURL); err == nil && !isS3URL(inputURL) {
		baseURL = u.Scheme + "://" + u.Host + u.Path
	}

	hash := sha256.Sum256([]byte(baseURL))
	return hex.EncodeToString(hash[:])[:16] // Use first 16 chars
}

// sanitizeURLForLog removes sensitive query parameters from URL for logging
func (d *S3Downloader) sanitizeURLForLog(inputURL string) string {
	if isS3URL(inputURL) {
		return inputURL
	}
	u, err := url.Parse(inputURL)
	if err != nil {
		return "[invalid URL]"
	}
	return u.Scheme + "://" + u.Host + u.Path + "?[signature hidden]"
}

// isS3URL checks if the URL is an s3:// URL
func isS3URL(inputURL string) bool {
	return strings.HasPrefix(inputURL, "s3://")
}

// parseS3URL parses an s3:// URL into bucket and key
func parseS3URL(inputURL string) (bucket, key string, err error) {
	if !isS3URL(inputURL) {
		return "", "", fmt.Errorf("not an S3 URL: %s", inputURL)
	}

	// Remove s3:// prefix
	path := strings.TrimPrefix(inputURL, "s3://")

	// Split into bucket and key
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid S3 URL format, expected s3://bucket/key: %s", inputURL)
	}

	return parts[0], parts[1], nil
}

// getS3ClientForBucket creates an S3 client configured for the bucket's region
func (d *S3Downloader) getS3ClientForBucket(ctx context.Context, bucket string) (*s3.Client, error) {
	// First, load the default config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create a client to detect the bucket region
	client := s3.NewFromConfig(cfg)

	// Get the actual bucket region
	region, err := s3manager.GetBucketRegion(ctx, client, bucket)
	if err != nil {
		// If we can't detect the region, return the default client
		if d.verbose {
			fmt.Printf("Warning: could not detect bucket region, using default: %v\n", err)
		}
		return client, nil
	}

	if d.verbose {
		fmt.Printf("Detected bucket region: %s\n", region)
	}

	// Reload config with the correct region
	cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config with region %s: %w", region, err)
	}

	return s3.NewFromConfig(cfg), nil
}

// isCacheValid checks if the cached file is still valid
func (d *S3Downloader) isCacheValid(ctx context.Context, inputURL, cachedFilePath, metadataPath string) bool {
	// Check if cached file exists
	if _, err := os.Stat(cachedFilePath); os.IsNotExist(err) {
		return false
	}

	// Check if metadata exists
	metadata, err := d.loadCacheMetadata(metadataPath)
	if err != nil {
		return false
	}

	// Verify URL hash matches
	if metadata.OriginalURL != d.hashURL(inputURL) {
		return false
	}

	// Get current ETag from remote
	var remoteETag string
	if isS3URL(inputURL) {
		remoteETag, err = d.getS3ETag(ctx, inputURL)
	} else {
		remoteETag, err = d.getPresignedURLETag(ctx, inputURL)
	}

	if err != nil {
		if d.verbose {
			fmt.Printf("Warning: failed to get remote ETag, will re-download: %v\n", err)
		}
		return false
	}

	// Compare ETags
	return metadata.ETag == remoteETag
}

// getS3ETag gets the ETag for an S3 object
func (d *S3Downloader) getS3ETag(ctx context.Context, inputURL string) (string, error) {
	bucket, key, err := parseS3URL(inputURL)
	if err != nil {
		return "", err
	}

	client, err := d.getS3ClientForBucket(ctx, bucket)
	if err != nil {
		return "", err
	}

	result, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get S3 object metadata: %w", err)
	}

	if result.ETag != nil {
		return *result.ETag, nil
	}
	return "", nil
}

// getPresignedURLETag gets the ETag for a presigned URL via HEAD request
func (d *S3Downloader) getPresignedURLETag(ctx context.Context, inputURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, inputURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HEAD request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HEAD request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HEAD request failed with status: %d", resp.StatusCode)
	}

	return resp.Header.Get("ETag"), nil
}

// downloadFromS3 downloads a file from S3 using the AWS SDK
func (d *S3Downloader) downloadFromS3(ctx context.Context, inputURL, destPath string) (string, error) {
	bucket, key, err := parseS3URL(inputURL)
	if err != nil {
		return "", err
	}

	client, err := d.getS3ClientForBucket(ctx, bucket)
	if err != nil {
		return "", err
	}

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return "", fmt.Errorf("failed to download from S3: %w", err)
	}
	defer func() {
		_ = result.Body.Close()
	}()

	// Create destination file
	file, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	// Copy content
	if _, err := io.Copy(file, result.Body); err != nil {
		_ = os.Remove(destPath) // Clean up partial file
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	var etag string
	if result.ETag != nil {
		etag = *result.ETag
	}

	return etag, nil
}

// downloadFromPresignedURL downloads a file from a presigned URL
func (d *S3Downloader) downloadFromPresignedURL(ctx context.Context, inputURL, destPath string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, inputURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GET request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute GET request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create destination file
	file, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	// Copy content
	if _, err := io.Copy(file, resp.Body); err != nil {
		_ = os.Remove(destPath) // Clean up partial file
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return resp.Header.Get("ETag"), nil
}

// loadCacheMetadata loads cache metadata from a JSON file
func (d *S3Downloader) loadCacheMetadata(path string) (*CacheMetadata, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata CacheMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// saveCacheMetadata saves cache metadata to a JSON file
func (d *S3Downloader) saveCacheMetadata(path string, metadata *CacheMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetDefaultCacheDir returns the default cache directory
func GetDefaultCacheDir() string {
	// Try user cache directory first
	if cacheDir, err := os.UserCacheDir(); err == nil {
		return filepath.Join(cacheDir, DefaultCacheSubdir)
	}

	// Fallback to home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".cache", DefaultCacheSubdir)
	}

	// Last resort: temp directory
	return filepath.Join(os.TempDir(), DefaultCacheSubdir)
}
