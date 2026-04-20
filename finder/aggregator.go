package finder

// GroupByHash groups hash results by their hash value and returns only
// groups where more than one file shares the same content.
func GroupByHash(results []HashResult) []DuplicateGroup {
	grouped := make(map[string][]HashResult)

	for _, r := range results {
		if r.Error != nil {
			continue
		}
		grouped[r.Hash] = append(grouped[r.Hash], r)
	}

	var groups []DuplicateGroup

	for hash, items := range grouped {
		if len(items) < 2 {
			continue
		}

		paths := make([]string, len(items))
		for i, item := range items {
			paths[i] = item.Path
		}

		size := items[0].Size
		waste := size * int64(len(items)-1)

		groups = append(groups, DuplicateGroup{
			Hash:       hash,
			Size:       size,
			Paths:      paths,
			TotalWaste: waste,
		})
	}

	return groups
}
