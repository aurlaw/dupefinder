package finder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"
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

// StartWorkers launches a pool of workerCount goroutines that read file candidates
// from jobs, hash each file, and send results to the returned channel.
// The results channel is closed automatically when all workers have finished.
func StartWorkers(ctx context.Context, jobs <-chan FileInfo, workerCount int) <-chan HashResult {
	results := make(chan HashResult, workerCount*2)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case f, ok := <-jobs:
					if !ok {
						return
					}
					hash, err := HashFile(f.Path)
					results <- HashResult{Path: f.Path, Size: f.Size, Hash: hash, Error: err}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
