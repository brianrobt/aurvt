#!/bin/bash

# Development build script for aurvt
# Usage: ./dev-build.sh [version] [commit] [build-type]

set -e

# Default values
VERSION=${1:-"0.0.0-dev"}
COMMIT=${2:-$(git rev-parse --short HEAD 2>/dev/null || echo "development")}
BUILD_TYPE=${3:-"dev"}

echo "=== aurvt Development Build ==="
echo "Version: $VERSION"
echo "Commit: $COMMIT"
echo "Build type: $BUILD_TYPE"
echo ""

# Build the binary
./build.sh "$VERSION" "$COMMIT" "$BUILD_TYPE"

echo ""
echo "=== Testing the build ==="
echo "Version info:"
./aurvt version

echo ""
echo "=== Testing on sample package ==="
./aurvt ../aur-pkgbuilds/python-conda-libmamba-solver

echo ""
echo "=== Build complete! ==="
echo "You can now test the binary with:"
echo "  ./aurvt version"
echo "  ./aurvt [package-directory]"