# dupefinder

A concurrent duplicate file finder CLI written in Go. Walks a directory tree, hashes files
using a configurable worker pool, and reports duplicates grouped by content hash.

----

Uses Go Releaser

```
brew install goreleaser

```

## Commands

```bash
# Build
go build -o dupefinder .

# Run
go run . find /path/to/scan
go run . find /path/to/scan --workers=8 --min-size=1KB --exclude=.git

# Test
go test ./...
go test ./finder/...          # single package
go test -run TestName ./...   # single test

# Vet
go vet ./...

# Example
go run . find ./testdata --min-size 1B

```


## Release

```
git tag v1.0.0
git push origin v1.0.0
```