# AUR Version Tool

A CLI tool to check if newer versions are available for AUR packages hosted on GitHub.

## Features

- Parse PKGBUILD files to extract package information
- Check GitHub releases for the latest version
- Compare current version with latest available
- Clean output
- **NEW**: Variable substitution support for PKGBUILD parsing
- **NEW**: Source array parsing with variable substitution
- **NEW**: Dual endpoint support (releases + tags)

## Installation

```bash
go install github.com/brianrobt/aurvt@latest
```

### Require in `go.mod`

```go
require github.com/brianrobt/aurvt v1.0.0
```

or

```go
import "github.com/brianrobt/aurvt@latest"
```

## Usage

```bash
aurvt <package-directory>
```

### Example

```bash
aurvt alist
```

Output:
```
Package: alist
Current version: 3.45.1
Repository URL: https://github.com/AlistGo/alist
Source URLs:
  [1] alist-3.45.1.tar.gz::https://github.com/AlistGo/alist/archive/refs/tags/3.45.1.tar.gz
Latest version: 3.46.2
🔄 New version available: 3.45.1 → 3.46.2
```

## Development

For local development and testing, see [DEVELOPMENT.md](./DEVELOPMENT.md).

### Quick Development Commands

```bash
# Build and test with development version
./dev-build.sh

# Test version bump
./bump-version.sh patch

# Test on multiple packages
make test-multiple
```

## Requirements

- Go 1.23+
- Internet connection for GitHub API access
- Valid PKGBUILD with `pkgname`, `pkgver`, and `url` fields

## Supported Repositories

Currently supports GitHub repositories only. The tool checks for:
- `url` field containing `github.com`
- Latest release via GitHub API (with fallback to tags)

## License

[GPL-3.0](./LICENSE)