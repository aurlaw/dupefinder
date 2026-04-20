package reporter

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aurlaw/dupefinder/finder"
	"github.com/dustin/go-humanize"
)

const (
	ansiRed   = "\033[31m"
	ansiDim   = "\033[2m"
	ansiReset = "\033[0m"
)

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func fmtElapsed(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.3fs", d.Seconds())
}

func colorize(s, code string) string {
	if isTerminal() {
		return code + s + ansiReset
	}
	return s
}

// PrintReport writes a human-readable duplicate report to stdout.
func PrintReport(r finder.Report) {
	if len(r.Groups) == 0 {
		fmt.Println("no duplicates found")
	} else {
		for i, group := range r.Groups {
			waste := humanize.Bytes(uint64(group.TotalWaste))
			fmt.Printf("group %d — %s each, %s wasted\n",
				i+1,
				humanize.Bytes(uint64(group.Size)),
				colorize(waste, ansiRed),
			)
			for _, path := range group.Paths {
				fmt.Printf("  %s\n", colorize(path, ansiDim))
			}
			fmt.Println()
		}
	}

	if len(r.Hardlinks) > 0 {
		fmt.Println("hardlinks (same inode, not wasted space):")
		for _, pair := range r.Hardlinks {
			fmt.Printf("  %s  ↔  %s\n", colorize(pair[0], ansiDim), colorize(pair[1], ansiDim))
		}
		fmt.Println()
	}

	printSummary(r)
}

func printSummary(r finder.Report) {
	fmt.Println("---")
	fmt.Printf("%-17s%s\n", "files scanned:", strconv.Itoa(r.TotalFiles))
	fmt.Printf("%-17s%s\n", "candidates:", strconv.Itoa(r.Candidates))
	fmt.Printf("%-17s%s\n", "duplicates:", strconv.Itoa(r.TotalDupes))
	fmt.Printf("%-17s%s\n", "hardlinks:", strconv.Itoa(len(r.Hardlinks)))
	fmt.Printf("%-17s%s\n", "reclaimable:", humanize.Bytes(uint64(r.WastedBytes)))
	fmt.Printf("%-17s%s\n", "elapsed:", fmtElapsed(r.ElapsedTime))
}
