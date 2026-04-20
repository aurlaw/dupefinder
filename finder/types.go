package finder

import "time"

// FileInfo represents a file candidate discovered during the walk.
type FileInfo struct {
	Path string
	Size int64
}

// HashResult is the output from the hasher for a single file.
type HashResult struct {
	Path  string
	Size  int64
	Hash  string
	Error error
}

// DuplicateGroup is a set of files with identical content.
type DuplicateGroup struct {
	Hash       string
	Size       int64
	Paths      []string
	TotalWaste int64
}

// Report is the final output passed to the reporter.
type Report struct {
	Groups      []DuplicateGroup
	TotalFiles  int
	Candidates  int
	TotalDupes  int
	WastedBytes int64
	ElapsedTime time.Duration
	Hardlinks   [][]string // pairs of paths that are hardlinks to the same inode
}
