package finder

// GroupBySize groups files by their size and returns only groups where
// at least two files share the same size. Files with unique sizes cannot
// be duplicates and are dropped without being hashed.
func GroupBySize(files []FileInfo) map[int64][]FileInfo {
	grouped := make(map[int64][]FileInfo)

	for _, f := range files {
		grouped[f.Size] = append(grouped[f.Size], f)
	}

	candidates := make(map[int64][]FileInfo)
	for size, group := range grouped {
		if len(group) > 1 {
			candidates[size] = group
		}
	}

	return candidates
}

// Flatten converts a size-grouped map back to a flat slice of FileInfo
// suitable for passing to the hasher.
func Flatten(grouped map[int64][]FileInfo) []FileInfo {
	var files []FileInfo
	for _, group := range grouped {
		files = append(files, group...)
	}
	return files
}
