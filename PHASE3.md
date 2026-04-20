# Phase 3 — concurrent worker pool

## Context
This is a duplicate file finder CLI. Phases 1 and 2 are complete:
- `finder.Walk` returns `[]FileInfo`
- `finder.GroupBySize` + `finder.Flatten` reduce candidates
- `finder.HashFile` hashes a single file sequentially
- `finder.GroupByHash` groups results into `[]DuplicateGroup`
- `reporter.PrintReport` outputs the final report

The sequential hashing loop in `cmd/find.go` now needs to be replaced
with a concurrent worker pool.

## Task
Implement `startWorkers` in `finder/hasher.go` per the stub and comments
already in place. Then update `cmd/find.go` to use it.

## Requirements
- Replace the sequential for loop in runFind with the channel-based approach
- Feed candidates into a buffered jobs channel from a separate goroutine
- Close the jobs channel when all candidates have been sent
- Range over the results channel to collect HashResults
- Use sync.WaitGroup to close results only when all workers are done
- Channel over mutex for the aggregation step
- Run go test -race ./... and fix any races before finishing

## Constraints
- Do not use sync.Mutex to protect any shared map
- startWorkers must accept ctx context.Context as first parameter
- Buffer size for jobs channel should be workerCount * 2
- All errors from HashFile flow through HashResult.Error, not panics