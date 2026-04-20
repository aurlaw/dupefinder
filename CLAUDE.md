# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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
```

## Architecture

`dupefinder` is a CLI tool (Cobra) that finds duplicate files by content hash. The `find` subcommand runs a 5-stage pipeline:

1. **Walk** (`finder/walker.go`) — Traverses the directory tree, filtering out excluded dirs, non-regular files, and files below `--min-size`.
2. **Size filter** (`finder/filter.go`) — Groups files by size and drops any file whose size is unique (cannot have duplicates), narrowing the hashing candidates.
3. **Hash** (`finder/hasher.go`) — Streams each candidate in 32 KB chunks to compute SHA-256. `startWorkers()` is a **stub** (panics); the `--workers` flag is wired up but concurrent hashing is not yet implemented.
4. **Aggregate** (`finder/aggregator.go`) — Groups hash results, keeps only groups with 2+ files, and computes wasted bytes per group.
5. **Report** (`reporter/text.go`) — Formats output using `go-humanize` for human-readable sizes and a summary of reclaimable space.

Orchestration lives in `cmd/find.go`. Types shared across packages are in `finder/types.go` (`FileInfo`, `HashResult`, `DuplicateGroup`, `Report`).

Test fixtures are in `testdata/` (a.txt, b.txt, c.txt, unique.txt).


