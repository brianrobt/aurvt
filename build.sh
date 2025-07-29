#!/bin/bash

# Build script for aurvt with version information

set -e

# Get version from git tag, or use default
VERSION=${1:-$(git describe --tags --abbrev=0 2>/dev/null || echo "0.1.0")}

# Get current commit hash
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "development")

# Get build date
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Building aurvt version $VERSION"
echo "Commit: $COMMIT"
echo "Date: $DATE"

# Build with version information
go build -ldflags "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.Date=$DATE" -o aurvt

echo "Build complete: aurvt"