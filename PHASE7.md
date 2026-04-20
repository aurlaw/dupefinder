# Phase 7 — JSON output and reporter polish

## Context
Phases 1–6 are complete. The pipeline is fully functional with concurrent
hashing, context cancellation, progress reporting, and symlink/hardlink
handling. `finder/types.go` has json struct tags on all types.
`reporter/json.go` has a `PrintJSON` stub to be implemented.

## Task
Implement the JSON reporter, switch between text and JSON output based on
the `--format` flag, and polish the text reporter output.

## Requirements

### 1. In `reporter/json.go`
- Implement `PrintJSON(r finder.Report) error`
- Marshal `Report` to JSON using `encoding/json`
- Use `json.MarshalIndent` with 2-space indentation for readability
- Write the output to stdout
- Return any marshalling or write error

### 2. In `reporter/text.go`
- Use `go-humanize` for ALL byte sizes — no manual formatting
- Polish the summary section to align columns using `fmt.Sprintf` padding:
  ```
  ---
  files scanned:    1024
  candidates:        312
  duplicates:         48
  hardlinks:           3
  reclaimable:     4.2 MB
  elapsed:         1.203s
  ```
- Add ANSI colour to waste totals (red) and paths (dim) only when stdout
  is a TTY — check with `isTerminal()` helper using `os.Stdout.Fd()`
- When stdout is not a TTY (e.g. piped), output plain text with no ANSI codes
- Print elapsed time using a human-readable format:
  - Under 1s: print as milliseconds e.g. `203ms`
  - Over 1s: print with 3 decimal places e.g. `1.203s`

### 3. In `cmd/find.go`
- After building the report, switch on the `--format` flag:
  ```go
  switch format {
  case "json":
      if err := reporter.PrintJSON(report); err != nil {
          return fmt.Errorf("writing json output: %w", err)
      }
  default:
      reporter.PrintReport(report)
  }
  ```
- If `--format` is not `text` or `json`, return a validation error before
  scanning begins — fail fast rather than scanning then failing

## Constraints
- ANSI codes only when stdout is a TTY — never when piped or redirected
- `PrintJSON` writes only valid JSON to stdout — no extra text or newlines
  before or after the JSON object
- `ElapsedMs` in the JSON output should be the duration in milliseconds
  as an integer, not a Go duration string — add a custom MarshalJSON or
  convert before marshalling
- Run `go test -race ./...` before finishing

## Acceptance

### Text output (default)
```bash
go run . find ./testdata --min-size 1B
```
Produces aligned, human-readable output with colour if run in a terminal.

### JSON output
```bash
go run . find ./testdata --min-size 1B --format json | jq .
```
Produces valid, well-structured JSON similar to:
```json
{
  "groups": [
    {
      "hash": "abc123...",
      "size": 6,
      "paths": ["testdata/a.txt", "testdata/b.txt"],
      "totalWaste": 6
    }
  ],
  "totalFiles": 4,
  "candidates": 2,
  "totalDupes": 2,
  "wastedBytes": 6,
  "elapsedMs": 1,
  "hardlinks": [
    ["testdata/a.txt", "testdata/a_hardlink.txt"]
  ]
}
```

### Invalid format flag
```bash
go run . find ./ --format xml
```
Returns an error before scanning:
```
Error: unsupported --format value "xml": must be "text" or "json"
```

### Piped output (no ANSI)
```bash
go run . find ./testdata --min-size 1B | cat
```
Produces plain text with no escape sequences.