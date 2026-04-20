package finder

import "time"

type FileInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type HashResult struct {
	Path  string `json:"path"`
	Size  int64  `json:"size"`
	Hash  string `json:"hash"`
	Error error  `json:"-"`
}

type DuplicateGroup struct {
	Hash       string   `json:"hash"`
	Size       int64    `json:"size"`
	Paths      []string `json:"paths"`
	TotalWaste int64    `json:"totalWaste"`
}

type Report struct {
	Groups      []DuplicateGroup `json:"groups"`
	TotalFiles  int              `json:"totalFiles"`
	Candidates  int              `json:"candidates"`
	TotalDupes  int              `json:"totalDupes"`
	WastedBytes int64            `json:"wastedBytes"`
	ElapsedTime time.Duration    `json:"elapsedMs"`
	Hardlinks   [][]string       `json:"hardlinks,omitempty"`
}
