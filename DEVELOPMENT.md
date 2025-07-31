# aurvt Development Guide

This guide explains how to build and test aurvt locally before creating releases.

## Quick Start

### 1. Development Build
```bash
# Build with development version
./dev-build.sh

# Or use make
make dev
```

### 2. Test Build with Custom Version
```bash
# Build with specific version
./build.sh "1.0.0-test" "abc123" "test"

# Or use make
make test-build
```

### 3. Version Bump Testing
```bash
# Bump patch version and test
./bump-version.sh patch

# Bump minor version and test
./bump-version.sh minor

# Bump major version and test
./bump-version.sh major

# Use custom version
./bump-version.sh patch "2.0.0-beta"
```

## Available Commands

### Build Commands
- `make build` - Build with version from git tag
- `make dev` - Development build with custom version
- `make test-build` - Test build with specific version
- `make quick` - Quick build without version injection
- `./dev-build.sh` - Development build with testing
- `./bump-version.sh [type] [version]` - Version bump and test

### Testing Commands
- `make test-package` - Test on sample package
- `make test-multiple` - Test on multiple packages
- `./aurvt version` - Show version information
- `./aurvt [package-directory]` - Test on specific package

### Utility Commands
- `make clean` - Remove build artifacts
- `make install` - Install to system
- `make uninstall` - Remove from system
- `make test` - Run Go tests
- `make help` - Show all available commands

## Development Workflow

### 1. Make Changes
Edit the code in `main.go` or other files.

### 2. Test Locally
```bash
# Build and test with development version
./dev-build.sh

# Or test specific version
./bump-version.sh patch
```

### 3. Test on Multiple Packages
```bash
make test-multiple
```

### 4. Commit and Push
```bash
git add .
git commit -m "feat: add new feature"
git push origin master
```

### 5. Release (Automatic)
The semantic-release setup will automatically:
- Analyze commits
- Generate changelog
- Create GitHub release
- Tag the release

## Version Management

### Current Version
The current version is managed in `package.json` and used by semantic-release.

### Version Bumping
- `patch` - Bug fixes (0.0.1)
- `minor` - New features (0.1.0)
- `major` - Breaking changes (1.0.0)

### Custom Versions
You can specify custom versions for testing:
```bash
./bump-version.sh patch "2.0.0-beta"
./build.sh "1.0.0-rc1" "abc123" "release"
```

## Testing Packages

### Sample Packages
The following packages are used for testing:
- `../aur-pkgbuilds/python-conda-libmamba-solver` - Standard package
- `../aur-pkgbuilds/rot8-git` - Git package
- `../aur-pkgbuilds/openmohaa-git` - Git package with variables

### Testing Commands
```bash
# Test single package
./aurvt ../aur-pkgbuilds/python-conda-libmamba-solver

# Test multiple packages
make test-multiple

# Test specific package
./aurvt ../aur-pkgbuilds/rot8-git
```

## Build Scripts

### build.sh
Main build script with version injection:
```bash
./build.sh [version] [commit] [build-type]
```

### dev-build.sh
Development build with automatic testing:
```bash
./dev-build.sh [version] [commit] [build-type]
```

### bump-version.sh
Version bump with testing:
```bash
./bump-version.sh [patch|minor|major] [custom-version]
```

## Troubleshooting

### Build Issues
- Ensure Go is installed and in PATH
- Check that all dependencies are available
- Verify git repository is properly initialized

### Testing Issues
- Ensure aur-pkgbuilds directory exists
- Check that sample packages are available
- Verify GitHub API access for version checking

### Version Issues
- Check package.json for current version
- Ensure git tags are properly set
- Verify semantic-release configuration

## Environment Variables

The build process uses these environment variables:
- `VERSION` - Version string (default: from git tag)
- `COMMIT` - Commit hash (default: from git)
- `DATE` - Build date (auto-generated)
- `BUILD_TYPE` - Build type (dev, test, release)

## Integration with CI/CD

The semantic-release setup automatically:
1. Analyzes commit messages
2. Determines version bump type
3. Generates changelog
4. Creates GitHub release
5. Tags the release

For local testing before release, use the development build commands above.