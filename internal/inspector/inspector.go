// Package inspector provides environment file inspection utilities,
// summarising key statistics and metadata about parsed .env entries.
package inspector

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

// Summary holds aggregate statistics for a parsed .env file.
type Summary struct {
	TotalKeys    int
	SecretKeys   int
	EmptyValues  int
	DuplicateKeys []string
	LongestKey   string
	Groups       []string
}

// Inspect parses the file at path and returns a Summary of its contents.
func Inspect(path string) (*Summary, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("inspect: %w", err)
	}
	return summarise(entries), nil
}

func summarise(entries []parser.Entry) *Summary {
	seen := make(map[string]int)
	groupSet := make(map[string]struct{})
	s := &Summary{}

	for _, e := range entries {
		seen[e.Key]++
		s.TotalKeys++

		if differ.IsSecret(e.Key) {
			s.SecretKeys++
		}
		if e.Value == "" {
			s.EmptyValues++
		}
		if len(e.Key) > len(s.LongestKey) {
			s.LongestKey = e.Key
		}
		if idx := strings.Index(e.Key, "_"); idx > 0 {
			groupSet[e.Key[:idx]] = struct{}{}
		}
	}

	for k, count := range seen {
		if count > 1 {
			s.DuplicateKeys = append(s.DuplicateKeys, k)
		}
	}
	for g := range groupSet {
		s.Groups = append(s.Groups, g)
	}
	return s
}
