#!/bin/bash

# Build script for aurvt with version information

set -e

# Default values
DEFAULT_VERSION="0.0.0-development"
DEFAULT_COMMIT="development"

# Get version from git tag, or use provided version, or use default
VERSION=${1:-$(git describe --tags --abbrev=0 2>/dev/null || echo "$DEFAULT_VERSION")}

# Get current commit hash
COMMIT=${2:-$(git rev-parse --short HEAD 2>/dev/null || echo "$DEFAULT_COMMIT")}

# Get build date
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build type (dev, release, test)
BUILD_TYPE=${3:-"dev"}

echo "Building aurvt version $VERSION"
echo "Commit: $COMMIT"
echo "Date: $DATE"
echo "Build type: $BUILD_TYPE"

# Build with version information
go build -ldflags "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.Date=$DATE" -o aurvt

echo "Build complete: aurvt"
echo ""
echo "To test the binary:"
echo "  ./aurvt version"
echo "  ./aurvt ../aur-pkgbuilds/python-conda-libmamba-solver"