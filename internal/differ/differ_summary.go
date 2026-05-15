package differ

import "github.com/user/envdiff/internal/parser"

// Summary holds aggregate statistics about a diff result.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Unchanged int
	Total    int
}

// Summarise computes aggregate counts from a slice of DiffEntry values.
func Summarise(entries []DiffEntry) Summary {
	var s Summary
	for _, e := range entries {
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusModified:
			s.Modified++
		case StatusUnchanged:
			s.Unchanged++
		}
	}
	s.Total = s.Added + s.Removed + s.Modified + s.Unchanged
	return s
}

// SummariseFiles parses two .env files and returns a Summary of their diff.
func SummariseFiles(baseFile, targetFile string) (Summary, error) {
	baseEntries, err := parser.Parse(baseFile)
	if err != nil {
		return Summary{}, err
	}
	targetEntries, err := parser.Parse(targetFile)
	if err != nil {
		return Summary{}, err
	}
	diffs := Diff(baseEntries, targetEntries)
	return Summarise(diffs), nil
}
