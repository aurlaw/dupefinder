package finder

import (
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
