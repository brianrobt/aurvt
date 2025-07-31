#!/bin/bash

# Version bump script for aurvt
# Usage: ./bump-version.sh [patch|minor|major] [new-version]

set -e

BUMP_TYPE=${1:-"patch"}
CUSTOM_VERSION=${2:-""}

# Get current version from package.json
CURRENT_VERSION=$(node -p "require('./package.json').version")

echo "Current version: $CURRENT_VERSION"

if [ -n "$CUSTOM_VERSION" ]; then
    NEW_VERSION="$CUSTOM_VERSION"
    echo "Using custom version: $NEW_VERSION"
else
    # Parse version components
    IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
    MAJOR=${VERSION_PARTS[0]}
    MINOR=${VERSION_PARTS[1]}
    PATCH=${VERSION_PARTS[2]//-*/}  # Remove development suffix

    case $BUMP_TYPE in
        "major")
            NEW_VERSION="$((MAJOR + 1)).0.0"
            ;;
        "minor")
            NEW_VERSION="$MAJOR.$((MINOR + 1)).0"
            ;;
        "patch")
            NEW_VERSION="$MAJOR.$MINOR.$((PATCH + 1))"
            ;;
        *)
            echo "Invalid bump type. Use: patch, minor, or major"
            exit 1
            ;;
    esac

    echo "Bumping $BUMP_TYPE version to: $NEW_VERSION"
fi

# Build with new version
echo ""
echo "Building with new version..."
./build.sh "$NEW_VERSION" "$(git rev-parse --short HEAD 2>/dev/null || echo 'development')" "test"

echo ""
echo "=== Testing new version ==="
./aurvt version

echo ""
echo "=== Testing on sample package ==="
./aurvt ../aur-pkgbuilds/python-conda-libmamba-solver

echo ""
echo "=== Version bump complete! ==="
echo "New version: $NEW_VERSION"
echo "You can test the binary with:"
echo "  ./aurvt version"
echo "  ./aurvt [package-directory]"