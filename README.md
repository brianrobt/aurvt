# AUR Version Tool

A CLI tool to check if newer versions are available for AUR packages hosted on GitHub.

## Features

- Parse PKGBUILD files to extract package information
- Check GitHub releases for the latest version
- Compare current version with latest available
- Clean output

## Installation

```bash
go build -o aurvt
```

## Usage

### Require in `go.mod`

```go
require github.com/brianrobt/aurvt v0.1.0
```

or

```go
import "github.com/brianrobt/aurvt@v0.1.0"
```

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
Latest version: 3.46.2
🔄 New version available: 3.45.1 → 3.46.2
```

## Requirements

- Go 1.23+
- Internet connection for GitHub API access
- Valid PKGBUILD with `pkgname`, `pkgver`, and `url` fields

## Supported Repositories

Currently supports GitHub repositories only. The tool checks for:
- `url` field containing `github.com`
- Latest release via GitHub API

## License

MIT