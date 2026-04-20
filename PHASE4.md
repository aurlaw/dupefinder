# Phase 4 — context cancellation and graceful shutdown

## Context
Phases 1–3 are complete. The worker pool in `finder/hasher.go` is implemented
and `cmd/find.go` uses it via channels. A placeholder `context.Background()`
is in place in `runFind` marked with a TODO comment.

## Task
Replace the placeholder context with proper signal handling so that Ctrl+C
causes a clean drain and exit rather than a hard kill.

## Requirements

1. In `cmd/find.go`:
   - Replace `context.Background()` with `signal.NotifyContext(context.Background(), os.Interrupt)`
   - Defer the cancel function
   - Import `os/signal` and `os`
   - After the results are collected, check if `ctx.Err() != nil` and print
     a cancellation message to stderr if so

2. In `finder/hasher.go`:
   - Each worker goroutine must respect `ctx.Done()` — use a select statement
     to choose between reading the next job and context cancellation
   - The file feeder goroutine in `cmd/find.go` must also check `ctx.Err()`
     before sending each job, stopping cleanly if cancelled

3. In `finder/hasher.go` `HashFile`:
   - Accept `ctx context.Context` as first parameter
   - Check `ctx.Err()` before opening the file — return early if cancelled

## Constraints
- Graceful shutdown must complete within 1-2 seconds of Ctrl+C
- No goroutine leaks — all goroutines must exit when context is cancelled
- Print a clear cancellation message to stderr, not stdout
- Do not use os.Exit — return an error from runFind instead
- Run `go test -race ./...` before finishing

## Acceptance
Running against a large directory and pressing Ctrl+C should produce:
  scan cancelled
  (partial summary if any results were collected)
Not a panic, not a hang.