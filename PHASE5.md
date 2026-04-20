# Phase 5 — progress reporting

## Context
Phases 1–4 are complete. The worker pool is running with context cancellation.
The `progress` package exists with a `Stats` type and a `StartReporter` stub.

## Task
Implement `StartReporter` in `progress/progress.go` and wire it into the
pipeline in `cmd/find.go`.

## Requirements

1. In `progress/progress.go`:
   - If `noProgress` is true, launch a goroutine that drains the stats
     channel and returns — do not print anything
   - Otherwise launch a goroutine with a `time.NewTicker(100 * time.Millisecond)`
   - Use `select` over three cases: incoming Stats, ticker tick, ctx.Done()
   - On each Stats message, accumulate into running totals (do not print yet)
   - On each ticker tick, reprint a single line to stderr using `\r` to
     overwrite the previous line
   - Format: `hashing... files: 42  bytes: 1.2 MB  dupes: 3`
   - Use `go-humanize` for byte formatting
   - Print a final newline to stderr when the goroutine exits cleanly

2. In `finder/hasher.go`:
   - Add a `stats chan<- progress.Stats` parameter to `startWorkers`
   - After each successful hash, send a Stats update with FilesHashed: 1
     and BytesHashed: file size
   - Send non-blocking — use a select with a default case so a slow
     progress reader never blocks a worker

3. In `cmd/find.go`:
   - Create a buffered stats channel (size 100) before starting workers
   - Call `progress.StartReporter` passing the channel
   - Close the stats channel after all results are collected
   - Pass `noProgress` flag through

## Constraints
- Progress output goes to stderr, results go to stdout — never mix them
- Non-blocking stats sends — workers must never block on progress reporting
- The stats channel must be closed after results are collected, not before
- Run `go test -race ./...` before finishing

## Acceptance
Running against a large directory shows a live updating line on stderr:
  hashing... files: 156  bytes: 45.2 MB  dupes: 4
Ctrl+C still exits cleanly with progress visible up to cancellation point.
`--no-progress` suppresses all progress output with no visible effect on
hashing performance.