#!/bin/sh

PACKAGE_VERSION=$(jq -r ".version" package.json)

echo "Building macOS distribution for stream-cli@${PACKAGE_VERSION}..."
pkg . --target=macos --out-path=./dist >/dev/null 2>&1

echo "Archiving build..."
tar -czf ./dist/stream-cli-${PACKAGE_VERSION}.tar.gz ./dist/getstream-cli

echo "Generating SHA-256..."
SHA=$(shasum -a 256 ./dist/stream-cli-${PACKAGE_VERSION}.tar.gz | cut -d' ' -f1)

echo "Uploading to GitHub..."
exec ./github.sh github_api_token=$GITHUB_TOKEN owner=GetStream repo=stream-cli tag=v${PACKAGE_VERSION} filename=./dist/stream-cli-${PACKAGE_VERSION}.tar.gz

echo "Build complete!"