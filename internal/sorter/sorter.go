// Package sorter provides functionality to sort .env file entries
// alphabetically or by key group, with optional section preservation.
package sorter

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// SortMode controls how entries are sorted.
type SortMode string

const (
	// SortAlpha sorts all keys alphabetically.
	SortAlpha SortMode = "alpha"
	// SortGroup sorts keys alphabetically within prefix groups (e.g. DB_, APP_).
	SortGroup SortMode = "group"
)

// Sort returns a new slice of entries sorted according to the given mode.
// The original slice is not modified.
func Sort(entries []parser.Entry, mode SortMode) []parser.Entry {
	result := make([]parser.Entry, len(entries))
	copy(result, entries)

	switch mode {
	case SortGroup:
		sortByGroup(result)
	default:
		sortAlpha(result)
	}

	return result
}

// sortAlpha sorts entries alphabetically by key.
func sortAlpha(entries []parser.Entry) {
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
}

// sortByGroup sorts entries by their prefix group first, then alphabetically
// within each group. A group is defined as the prefix before the first '_'.
func sortByGroup(entries []parser.Entry) {
	sort.SliceStable(entries, func(i, j int) bool {
		gi := groupOf(entries[i].Key)
		gj := groupOf(entries[j].Key)
		if gi != gj {
			return gi < gj
		}
		return entries[i].Key < entries[j].Key
	})
}

// groupOf returns the prefix group of a key (part before the first '_').
// If there is no '_', the full key is used as the group.
func groupOf(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
