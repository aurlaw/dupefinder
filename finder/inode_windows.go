//go:build windows

package finder

import "os"

// InodeKey is not reliably supported on Windows.
// Hardlink detection is skipped on this platform.
type InodeKey struct {
	Dev uint64
	Ino uint64
}

// GetInodeKey always returns an empty InodeKey on Windows.
func GetInodeKey(info os.FileInfo) (InodeKey, error) {
	return InodeKey{}, nil
}
