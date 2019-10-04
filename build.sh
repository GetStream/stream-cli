#!/bin/sh

PACKAGE_VERSION=$(jq -r ".version" package.json)

echo "Building macOS distribution for stream-cli@${PACKAGE_VERSION}..."
pkg . --target=macos --out-path=./dist >/dev/null 2>&1

echo "Archiving build..."
tar -czf ./dist/stream-cli-${PACKAGE_VERSION}.tar.gz ./dist/getstream-cli

echo "Generating SHA-256..."
SHA=$(shasum -a 256 ./dist/stream-cli-${PACKAGE_VERSION}.tar.gz | cut -d' ' -f1)

echo "Build complete!"