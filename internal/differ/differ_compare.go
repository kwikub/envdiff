package differ

import "github.com/subtlepseudonym/envdiff/internal/parser"

// CompareResult holds the comparison of a single key across two env files.
type CompareResult struct {
	Key      string
	Left     string
	Right    string
	Status   string // "match", "mismatch", "left_only", "right_only"
}

// Compare performs a key-by-key comparison of two entry slices and returns
// a CompareResult for every key seen in either side.
func Compare(left, right []parser.Entry) []CompareResult {
	lm := toMap(left)
	rm := toMap(right)

	seen := make(map[string]struct{})
	var results []CompareResult

	for k, lv := range lm {
		seen[k] = struct{}{}
		if rv, ok := rm[k]; ok {
			status := "match"
			if lv != rv {
				status = "mismatch"
			}
			results = append(results, CompareResult{Key: k, Left: lv, Right: rv, Status: status})
		} else {
			results = append(results, CompareResult{Key: k, Left: lv, Right: "", Status: "left_only"})
		}
	}

	for k, rv := range rm {
		if _, ok := seen[k]; !ok {
			results = append(results, CompareResult{Key: k, Left: "", Right: rv, Status: "right_only"})
		}
	}

	return results
}
