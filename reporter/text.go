package reporter

import (
	"fmt"

	"github.com/aurlaw/dupefinder/finder"
	"github.com/dustin/go-humanize"
)

// PrintReport writes a human-readable duplicate report to stdout.
func PrintReport(r finder.Report) {
	if len(r.Groups) == 0 {
		fmt.Println("no duplicates found")
		printSummary(r)
		return
	}

	for i, group := range r.Groups {
		fmt.Printf("group %d — %s each, %s wasted\n",
			i+1,
			humanize.Bytes(uint64(group.Size)),
			humanize.Bytes(uint64(group.TotalWaste)),
		)
		for _, path := range group.Paths {
			fmt.Printf("  %s\n", path)
		}
		fmt.Println()
	}

	if len(r.Hardlinks) > 0 {
		fmt.Println("hardlinks (same inode, not wasted space):")
		for _, pair := range r.Hardlinks {
			fmt.Printf("  %s  ↔  %s\n", pair[0], pair[1])
		}
		fmt.Println()
	}

	printSummary(r)
}

func printSummary(r finder.Report) {
	fmt.Println("---")
	fmt.Printf("files scanned:   %d\n", r.TotalFiles)
	fmt.Printf("candidates:      %d\n", r.Candidates)
	fmt.Printf("duplicates:      %d\n", r.TotalDupes)
	fmt.Printf("hardlinks:       %d\n", len(r.Hardlinks))
	fmt.Printf("reclaimable:     %s\n", humanize.Bytes(uint64(r.WastedBytes)))
	fmt.Printf("elapsed:         %s\n", r.ElapsedTime)
}
