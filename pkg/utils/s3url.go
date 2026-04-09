package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// awsRegionPattern matches typical AWS region identifiers (e.g. us-east-1, eu-central-1).
var awsRegionPattern = regexp.MustCompile(`^[a-z]{2}(-gov)?-[a-z]+-\d+$`)

// regionFromAmzCredential returns the AWS region from the X-Amz-Credential query
// parameter (format: accessKey/date/region/service/aws4_request).
func regionFromAmzCredential(u *url.URL) string {
	v := u.Query().Get("X-Amz-Credential")
	if v == "" {
		return ""
	}
	decoded, err := url.QueryUnescape(v)
	if err != nil {
		decoded = v
	}
	parts := strings.Split(decoded, "/")
	if len(parts) < 3 {
		return ""
	}
	r := parts[2]
	if !awsRegionPattern.MatchString(r) {
		return ""
	}
	return r
}

// S3BucketAndRegionFromUploadURL derives the S3 bucket and region from a presigned PUT URL.
// It handles virtual-hosted–style hosts (bucket.s3.region.amazonaws.com, bucket.s3.amazonaws.com)
// and falls back to the region embedded in X-Amz-Credential when the host is non-standard
// (e.g. S3 Transfer Acceleration).
func S3BucketAndRegionFromUploadURL(raw string) (bucket, region string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", fmt.Errorf("parse upload URL: %w", err)
	}
	host := strings.ToLower(u.Hostname())
	credRegion := regionFromAmzCredential(u)
	parts := strings.Split(host, ".")
	n := len(parts)

	if n >= 4 && parts[n-2] == "amazonaws" && (parts[n-1] == "com" || parts[n-1] == "cn") {
		switch {
		case n == 4 && parts[1] == "s3":
			// e.g. bucket.s3.amazonaws.com → us-east-1 for the global endpoint
			return parts[0], "us-east-1", nil
		case n == 5 && parts[1] == "s3":
			// e.g. bucket.s3.eu-central-1.amazonaws.com
			return parts[0], parts[2], nil
		case n == 6 && parts[1] == "s3" && parts[2] == "dualstack":
			return parts[0], parts[3], nil
		}
	}

	// Path-style: s3.region.amazonaws.com/bucket/key
	if n == 4 && parts[0] == "s3" && parts[2] == "amazonaws" && awsRegionPattern.MatchString(parts[1]) {
		reg := parts[1]
		path := strings.Trim(strings.TrimPrefix(u.Path, "/"), "/")
		if path == "" {
			return "", "", fmt.Errorf("path-style S3 URL missing object key path")
		}
		slash := strings.IndexByte(path, '/')
		if slash <= 0 {
			return "", "", fmt.Errorf("path-style S3 URL missing bucket in path")
		}
		return path[:slash], reg, nil
	}

	// Fallback: first label before ".s3" is the bucket (virtual-hosted variants like accelerate)
	if i := strings.Index(host, ".s3"); i > 0 {
		b := host[:i]
		if b != "" {
			r := credRegion
			if r == "" && n == 4 && parts[1] == "s3" {
				r = "us-east-1"
			}
			if r != "" {
				return b, r, nil
			}
		}
	}

	if credRegion != "" {
		return "", "", fmt.Errorf("could not derive S3 bucket from upload URL host %q (region from credential: %s)", host, credRegion)
	}
	return "", "", fmt.Errorf("could not derive S3 bucket and region from upload URL host %q", host)
}
