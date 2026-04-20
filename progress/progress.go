package progress

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

// Stats represents a snapshot of hashing progress sent from workers.
type Stats struct {
	FilesHashed int64
	BytesHashed int64
	DupesFound  int64
}

// StartReporter launches a goroutine that reads from the stats channel and
// prints a live progress line to stderr every 100ms.
// The returned channel closes when the goroutine exits, allowing the caller
// to wait for the final newline before printing the report.
func StartReporter(ctx context.Context, stats <-chan Stats, noProgress bool) <-chan struct{} {
	done := make(chan struct{})

	if noProgress {
		go func() {
			defer close(done)
			for range stats {
			}
		}()
		return done
	}

	go func() {
		defer close(done)
		defer fmt.Fprint(os.Stderr, "\n")

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		var files, bytes, dupes int64
		for {
			select {
			case s, ok := <-stats:
				if !ok {
					fmt.Fprintf(os.Stderr, "\rhashing... files: %d  bytes: %s  dupes: %d",
						files, humanize.Bytes(uint64(bytes)), dupes)
					return
				}
				files += s.FilesHashed
				bytes += s.BytesHashed
				dupes += s.DupesFound
			case <-ticker.C:
				fmt.Fprintf(os.Stderr, "\rhashing... files: %d  bytes: %s  dupes: %d",
					files, humanize.Bytes(uint64(bytes)), dupes)
			case <-ctx.Done():
				return
			}
		}
	}()

	return done
}
