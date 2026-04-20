package reporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aurlaw/dupefinder/finder"
)

// jsonReport is a serialization-only struct that converts ElapsedTime to integer milliseconds.
type jsonReport struct {
	Groups      []finder.DuplicateGroup `json:"groups"`
	TotalFiles  int                     `json:"totalFiles"`
	Candidates  int                     `json:"candidates"`
	TotalDupes  int                     `json:"totalDupes"`
	WastedBytes int64                   `json:"wastedBytes"`
	ElapsedMs   int64                   `json:"elapsedMs"`
	Hardlinks   [][]string              `json:"hardlinks,omitempty"`
}

// PrintJSON writes a machine-readable JSON report to stdout.
func PrintJSON(r finder.Report) error {
	out := jsonReport{
		Groups:      r.Groups,
		TotalFiles:  r.TotalFiles,
		Candidates:  r.Candidates,
		TotalDupes:  r.TotalDupes,
		WastedBytes: r.WastedBytes,
		ElapsedMs:   r.ElapsedTime.Milliseconds(),
		Hardlinks:   r.Hardlinks,
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stdout, "%s\n", data)
	return err
}
