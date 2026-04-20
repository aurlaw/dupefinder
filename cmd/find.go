package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/aurlaw/dupefinder/finder"
	"github.com/aurlaw/dupefinder/reporter"
)

var (
	workers    int
	minSize    string
	excludes   []string
	format     string
	noProgress bool
)

var findCmd = &cobra.Command{
	Use:   "find [path]",
	Short: "Find duplicate files in the given directory",
	Args:  cobra.ExactArgs(1),
	RunE:  runFind,
}

func init() {
	findCmd.Flags().IntVar(&workers, "workers", 8, "Number of concurrent hashing goroutines")
	findCmd.Flags().StringVar(&minSize, "min-size", "1KB", "Skip files smaller than this size")
	findCmd.Flags().StringArrayVar(&excludes, "exclude", []string{".git"}, "Directory names to skip")
	findCmd.Flags().StringVar(&format, "format", "text", "Output format: text or json")
	findCmd.Flags().BoolVar(&noProgress, "no-progress", false, "Suppress progress output")

	rootCmd.AddCommand(findCmd)
}

func runFind(cmd *cobra.Command, args []string) error {
	root := args[0]
	start := time.Now()

	// Parse min-size flag into bytes
	minBytes, err := humanize.ParseBytes(minSize)
	if err != nil {
		return fmt.Errorf("invalid --min-size value %q: %w", minSize, err)
	}

	// Step 1: walk the directory tree
	fmt.Fprintf(os.Stderr, "walking %s...\n", root)
	files, err := finder.Walk(root, excludes, int64(minBytes))
	if err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}
	fmt.Fprintf(os.Stderr, "found %d files\n", len(files))

	// Step 2: pre-filter by size — drop files with unique sizes
	sizeGroups := finder.GroupBySize(files)
	candidates := finder.Flatten(sizeGroups)
	fmt.Fprintf(os.Stderr, "%d candidates after size filter\n", len(candidates))

	// Step 3: hash candidates only
	fmt.Fprintf(os.Stderr, "hashing files...\n")
	results := make([]finder.HashResult, 0, len(candidates))
	for _, f := range candidates {
		hash, err := finder.HashFile(f.Path)
		results = append(results, finder.HashResult{
			Path:  f.Path,
			Size:  f.Size,
			Hash:  hash,
			Error: err,
		})
	}

	// Step 4: group by hash, find duplicates
	groups := finder.GroupByHash(results)

	// Step 5: build report
	var totalWaste int64
	var totalDupes int
	for _, g := range groups {
		totalWaste += g.TotalWaste
		totalDupes += len(g.Paths)
	}

	report := finder.Report{
		Groups:      groups,
		TotalFiles:  len(files),
		Candidates:  len(candidates),
		TotalDupes:  totalDupes,
		WastedBytes: totalWaste,
		ElapsedTime: time.Since(start),
	}

	// Step 5: print report
	reporter.PrintReport(report)

	return nil
}
