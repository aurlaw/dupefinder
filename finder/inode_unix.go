//go:build !windows

package finder

import (
	"os"
	"syscall"
)

// InodeKey uniquely identifies a file by device and inode number.
// Two paths with the same InodeKey are hardlinks to the same data.
type InodeKey struct {
	Dev uint64
	Ino uint64
}

// GetInodeKey extracts the device and inode number from a file's metadata.
func GetInodeKey(info os.FileInfo) (InodeKey, error) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return InodeKey{}, nil
	}
	return InodeKey{
		Dev: uint64(stat.Dev),
		Ino: uint64(stat.Ino),
	}, nil
}
