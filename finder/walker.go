package finder

import (
	"io/fs"
	"path/filepath"
)

// Walk traverses root recursively, returning all files that meet the minimum
// size threshold. Directories matching any name in excludes are skipped entirely.
func Walk(root string, excludes []string, minSize int64) ([]FileInfo, error) {
	var files []FileInfo

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

		// Skip non-regular files (symlinks, pipes, devices)
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

		files = append(files, FileInfo{
			Path: path,
			Size: info.Size(),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
