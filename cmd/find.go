package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/aurlaw/dupefinder/finder"
	"github.com/aurlaw/dupefinder/progress"
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

	if format != "text" && format != "json" {
		return fmt.Errorf("unsupported --format value %q: must be \"text\" or \"json\"", format)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Parse min-size flag into bytes
	minBytes, err := humanize.ParseBytes(minSize)
	if err != nil {
		return fmt.Errorf("invalid --min-size value %q: %w", minSize, err)
	}

	// Step 1: walk the directory tree
	fmt.Fprintf(os.Stderr, "walking %s...\n", root)
	files, hardlinks, err := finder.Walk(root, excludes, int64(minBytes))
	if err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}
	fmt.Fprintf(os.Stderr, "found %d files\n", len(files))

	// Step 2: pre-filter by size — drop files with unique sizes
	sizeGroups := finder.GroupBySize(files)
	candidates := finder.Flatten(sizeGroups)
	fmt.Fprintf(os.Stderr, "%d candidates after size filter\n", len(candidates))

	// Step 3: hash candidates concurrently
	statsCh := make(chan progress.Stats, 100)
	reporterDone := progress.StartReporter(ctx, statsCh, noProgress)

	jobs := make(chan finder.FileInfo, workers*2)
	go func() {
		defer close(jobs)
		for _, f := range candidates {
			select {
			case jobs <- f:
			case <-ctx.Done():
				return
			}
		}
	}()

	resultsCh := finder.StartWorkers(ctx, jobs, workers, statsCh)

	results := make([]finder.HashResult, 0, len(candidates))
	for r := range resultsCh {
		results = append(results, r)
	}

	close(statsCh)
	<-reporterDone

	if ctx.Err() != nil {
		fmt.Fprintln(os.Stderr, "scan cancelled")
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
		Hardlinks:   hardlinks,
		TotalFiles:  len(files),
		Candidates:  len(candidates),
		TotalDupes:  totalDupes,
		WastedBytes: totalWaste,
		ElapsedTime: time.Since(start),
	}

	// Step 5: print report
	switch format {
	case "json":
		if err := reporter.PrintJSON(report); err != nil {
			return fmt.Errorf("writing json output: %w", err)
		}
	default:
		reporter.PrintReport(report)
	}

	return nil
}
