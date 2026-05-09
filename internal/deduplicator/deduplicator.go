// Package deduplicator removes duplicate keys from a slice of env entries,
// applying a configurable strategy to decide which occurrence to keep.
package deduplicator

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Strategy controls which duplicate entry is retained.
type Strategy string

const (
	// StrategyFirst keeps the first occurrence of a duplicate key.
	StrategyFirst Strategy = "first"
	// StrategyLast keeps the last occurrence of a duplicate key.
	StrategyLast Strategy = "last"
)

// Result holds the deduplicated entries and a log of removed duplicates.
type Result struct {
	Entries  []parser.Entry
	Removed  []Duplicate
}

// Duplicate records a key that appeared more than once and the line that was dropped.
type Duplicate struct {
	Key     string
	Line    int
	Value   string
}

// Deduplicate removes duplicate keys from entries according to the given strategy.
// It returns an error if strategy is unrecognised.
func Deduplicate(entries []parser.Entry, strategy Strategy) (Result, error) {
	if strategy != StrategyFirst && strategy != StrategyLast {
		return Result{}, fmt.Errorf("deduplicator: unknown strategy %q", strategy)
	}

	seen := make(map[string]int) // key -> index in out slice
	out := make([]parser.Entry, 0, len(entries))
	var removed []Duplicate

	for _, e := range entries {
		if idx, exists := seen[e.Key]; exists {
			if strategy == StrategyLast {
				// Record the entry currently in out as removed, replace it.
				prev := out[idx]
				removed = append(removed, Duplicate{Key: prev.Key, Line: prev.Line, Value: prev.Value})
				out[idx] = e
			} else {
				// StrategyFirst: discard the new entry.
				removed = append(removed, Duplicate{Key: e.Key, Line: e.Line, Value: e.Value})
			}
		} else {
			seen[e.Key] = len(out)
			out = append(out, e)
		}
	}

	return Result{Entries: out, Removed: removed}, nil
}
