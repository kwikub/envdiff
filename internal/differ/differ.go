package differ

import "github.com/user/envdiff/internal/parser"

// DiffType represents the type of difference between two env files.
type DiffType string

const (
	Added    DiffType = "added"
	Removed  DiffType = "removed"
	Modified DiffType = "modified"
	Unchanged DiffType = "unchanged"
)

// DiffEntry represents a single key difference between two env files.
type DiffEntry struct {
	Key      string
	BaseVal  string
	OtherVal string
	Type     DiffType
}

// Diff compares two sets of parsed env entries and returns the differences.
func Diff(base, other []parser.Entry) []DiffEntry {
	baseMap := toMap(base)
	otherMap := toMap(other)

	var results []DiffEntry

	// Check for removed or modified keys
	for _, entry := range base {
		if otherVal, ok := otherMap[entry.Key]; ok {
			if otherVal != entry.Value {
				results = append(results, DiffEntry{
					Key:      entry.Key,
					BaseVal:  entry.Value,
					OtherVal: otherVal,
					Type:     Modified,
				})
			} else {
				results = append(results, DiffEntry{
					Key:      entry.Key,
					BaseVal:  entry.Value,
					OtherVal: otherVal,
					Type:     Unchanged,
				})
			}
		} else {
			results = append(results, DiffEntry{
				Key:     entry.Key,
				BaseVal: entry.Value,
				Type:    Removed,
			})
		}
	}

	// Check for added keys
	for _, entry := range other {
		if _, ok := baseMap[entry.Key]; !ok {
			results = append(results, DiffEntry{
				Key:      entry.Key,
				OtherVal: entry.Value,
				Type:     Added,
			})
		}
	}

	return results
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
