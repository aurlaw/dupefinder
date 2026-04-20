package finder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const chunkSize = 32 * 1024 // 32KB read buffer

// HashFile reads the file at path and returns its SHA-256 hash as a hex string.
// The file is streamed in chunks to avoid loading large files into memory.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	buf := make([]byte, chunkSize)

	if _, err := io.CopyBuffer(h, f, buf); err != nil {
		return "", fmt.Errorf("hashing %s: %w", path, err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// startWorkers launches a pool of workerCount goroutines that read file candidates
// from jobs, hash each file, and send results to the returned channel.
// The results channel is closed automatically when all workers have finished.
func startWorkers(ctx context.Context, jobs <-chan FileInfo, workerCount int) <-chan HashResult {
	// TODO: implement worker pool
	// - create results channel (buffered, size = workerCount * 2)
	// - spawn workerCount goroutines, each ranging over jobs
	// - use sync.WaitGroup to close results when all workers done
	// - each worker calls HashFile and sends a HashResult
	// - respect ctx.Done() for cancellation
	panic("not implemented")
}
