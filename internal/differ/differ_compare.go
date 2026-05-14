package differ

// CompareResult holds the result of comparing two env files side by side.
type CompareResult struct {
	Key    string
	Left   string
	Right  string
	Status string // "match", "mismatch", "left_only", "right_only"
}

// Compare performs a side-by-side comparison of two entry slices and returns
// a list of CompareResults describing each key's status across both sides.
func Compare(left, right []Entry) []CompareResult {
	leftMap := toMap(left)
	rightMap := toMap(right)

	seen := make(map[string]bool)
	var results []CompareResult

	for _, e := range left {
		if seen[e.Key] {
			continue
		}
		seen[e.Key] = true

		rv, ok := rightMap[e.Key]
		if !ok {
			results = append(results, CompareResult{
				Key:    e.Key,
				Left:   e.Value,
				Right:  "",
				Status: "left_only",
			})
			continue
		}
		status := "match"
		if e.Value != rv {
			status = "mismatch"
		}
		results = append(results, CompareResult{
			Key:    e.Key,
			Left:   e.Value,
			Right:  rv,
			Status: status,
		})
	}

	for _, e := range right {
		if seen[e.Key] {
			continue
		}
		seen[e.Key] = true
		if _, ok := leftMap[e.Key]; !ok {
			results = append(results, CompareResult{
				Key:    e.Key,
				Left:   "",
				Right:  e.Value,
				Status: "right_only",
			})
		}
	}

	return results
}
