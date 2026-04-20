# Phase 6 — symlink and hardlink handling

## Context
Phases 1–5 are complete. The pipeline walks, filters, hashes concurrently,
cancels cleanly, and reports progress. `finder/inode_unix.go` and
`finder/inode_windows.go` provide a platform-abstracted `GetInodeKey`
function. `Report` has a `Hardlinks [][]string` field.

## Task
Update `finder/walker.go` to detect and skip symlinks, and to identify
hardlinks using inode tracking. Update the reporter to display hardlinks
separately from content duplicates.

## Requirements

### 1. In `finder/walker.go`
- Use `os.Lstat` via `d.Info()` — `WalkDir` already uses Lstat internally
  so `d.Type()` correctly reflects symlinks without following them
- Skip symlinks by checking `d.Type()&fs.ModeSymlink != 0` — return nil
  to skip without error
- Add a `seenInodes map[InodeKey]string` local variable to `Walk`
- For each regular file, call `GetInodeKey` with the file info
- If the inode key is already in `seenInodes`, record both paths as a
  hardlink pair and skip the file (do not add to candidates)
- If the inode key is zero value (Windows), skip inode tracking entirely
- Return hardlink pairs as `[][]string` alongside `[]FileInfo`
- Update the `Walk` signature to:
  `func Walk(root string, excludes []string, minSize int64) ([]FileInfo, [][]string, error)`

### 2. In `cmd/find.go`
- Update the call to `finder.Walk` to capture the hardlinks return value
- Pass hardlinks into the `Report`

### 3. In `reporter/text.go`
- After printing duplicate groups, print a hardlinks section if any exist
- Format:
  ```
  hardlinks (same inode, not wasted space):
    path/to/file-a  ↔  path/to/file-b
  ```
- Print hardlink count in the summary line

## Constraints
- Symlinks must be skipped silently — no error, no output
- Hardlinks must not appear in duplicate groups
- The zero-value InodeKey check guards Windows — do not use build tags
  in walker.go itself, only in the inode_*.go files
- Run `go test -race ./...` before finishing

## Acceptance
Create test hardlinks and symlinks to verify:

```bash
ln testdata/a.txt testdata/a_hardlink.txt
ln -s testdata/a.txt testdata/a_symlink.txt
go run . find ./testdata --min-size 1B
```

Running the finder should produce output similar to:

```
group 1 — 6 B each, 6 B wasted
  testdata/a.txt
  testdata/b.txt

hardlinks (same inode, not wasted space):
  testdata/a.txt  ↔  testdata/a_hardlink.txt
---
files scanned:   4
candidates:      2
duplicates:      2
hardlinks:       1
reclaimable:     6 B
elapsed:         ...
```

Specifically:
- a.txt and b.txt reported as content duplicates
- a.txt and a_hardlink.txt reported as hardlinks
- a_symlink.txt silently skipped with no error