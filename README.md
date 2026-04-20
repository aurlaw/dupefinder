# dupefinder

A concurrent duplicate file finder CLI written in Go. Walks a directory tree, hashes files
using a configurable worker pool, and reports duplicates grouped by content hash.

---

## Installation

### Download (macOS, Linux, Windows)

Download the latest binary for your platform from the
[releases page](https://github.com/aurlaw/dupefinder/releases/latest):

| Platform | Archive |
|---|---|
| macOS (Apple Silicon) | `dupefinder_darwin_arm64.tar.gz` |
| macOS (Intel) | `dupefinder_darwin_amd64.tar.gz` |
| Linux (amd64) | `dupefinder_linux_amd64.tar.gz` |
| Windows (amd64) | `dupefinder_windows_amd64.zip` |

Extract and move the binary to a directory on your PATH:

```bash
# macOS/Linux
tar -xzf dupefinder_darwin_arm64.tar.gz
mv dupefinder /usr/local/bin/
```

### Go install
```bash
go install github.com/aurlaw/dupefinder@latest
```

---

## Usage

```bash
# Scan a directory with default settings
dupefinder find /path/to/scan

# Tune the worker pool and minimum file size
dupefinder find /path/to/scan --workers 8 --min-size 1MB

# Exclude additional directories
dupefinder find /path/to/scan --exclude .git --exclude node_modules

# Output as JSON
dupefinder find /path/to/scan --format json

# Output as JSON and pipe to jq
dupefinder find /path/to/scan --format json | jq .

# Suppress progress output
dupefinder find /path/to/scan --no-progress
```

## Flags

| Flag | Default | Description |
|---|---|---|
| `--workers` | `8` | Number of concurrent hashing goroutines |
| `--min-size` | `1KB` | Skip files smaller than this size |
| `--exclude` | `.git` | Directory names to skip (repeatable) |
| `--format` | `text` | Output format: `text` or `json` |
| `--no-progress` | `false` | Suppress progress output on stderr |

---

## Development

### Prerequisites
- Go 1.22+
- [GoReleaser](https://goreleaser.com) (for releases only)

```bash
brew install goreleaser
```

### Commands

```bash
# Build
go build -o dupefinder .

# Run
go run . find /path/to/scan

# Test
go test ./...
go test ./finder/...          # single package
go test -run TestName ./...   # single test

# Race detector
go test -race ./...

# Vet
go vet ./...
```

### Releasing

Releases are automated via GoReleaser and GitHub Actions. Pushing a version
tag triggers a full build for all platforms and publishes a GitHub Release.

```bash
git tag v1.0.0
git push origin v1.0.0
```

Binaries are stamped with the tag version at build time. Local builds without
GoReleaser report `dev` as the version:

```bash
dupefinder --version
# dupefinder version dev
```