package finder

import (
	"io/fs"
	"path/filepath"
)

// Walk traverses root recursively, returning all files that meet the minimum
// size threshold and a list of hardlink pairs detected by inode tracking.
// Directories matching any name in excludes are skipped entirely.
// Symlinks are silently skipped. Hardlinked files appear only once in the
// returned FileInfo slice; the pair is recorded in hardlinks instead.
func Walk(root string, excludes []string, minSize int64) ([]FileInfo, [][]string, error) {
	var files []FileInfo
	var hardlinks [][]string
	seenInodes := make(map[InodeKey]string)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if d.IsDir() {
			for _, ex := range excludes {
				if d.Name() == ex {
					return fs.SkipDir
				}
			}
			return nil
		}

		// Skip symlinks explicitly
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}

		// Skip non-regular files (pipes, devices, etc.)
		if !d.Type().IsRegular() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		// Skip files below minimum size
		if info.Size() < minSize {
			return nil
		}

		// Hardlink detection via inode key
		key, err := GetInodeKey(info)
		if err != nil {
			return err
		}
		if key != (InodeKey{}) { // zero value means Windows — skip tracking
			if first, seen := seenInodes[key]; seen {
				hardlinks = append(hardlinks, []string{first, path})
				return nil // already counted under the first path
			}
			seenInodes[key] = path
		}

		files = append(files, FileInfo{
			Path: path,
			Size: info.Size(),
		})

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return files, hardlinks, nil
}
