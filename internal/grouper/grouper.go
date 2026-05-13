// Package grouper organises env entries into named groups based on key prefixes.
package grouper

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Group holds a named collection of env entries sharing a common prefix.
type Group struct {
	Name    string
	Entries []parser.Entry
}

// Options controls how grouping is performed.
type Options struct {
	// Separator is the delimiter used to split the prefix from the rest of the
	// key. Defaults to "_".
	Separator string
	// IncludeUngrouped controls whether keys with no prefix are returned in a
	// synthetic group named "(ungrouped)".
	IncludeUngrouped bool
}

// GroupBy partitions entries into groups derived from each key's prefix segment.
// Keys that contain no separator are placed in the "(ungrouped)" bucket when
// opts.IncludeUngrouped is true, otherwise they are silently dropped.
func GroupBy(entries []parser.Entry, opts Options) []Group {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	buckets := make(map[string][]parser.Entry)

	for _, e := range entries {
		idx := strings.Index(e.Key, opts.Separator)
		var prefix string
		if idx <= 0 {
			prefix = "(ungrouped)"
		} else {
			prefix = e.Key[:idx]
		}

		if prefix == "(ungrouped)" && !opts.IncludeUngrouped {
			continue
		}
		buckets[prefix] = append(buckets[prefix], e)
	}

	names := make([]string, 0, len(buckets))
	for k := range buckets {
		names = append(names, k)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, name := range names {
		groups = append(groups, Group{Name: name, Entries: buckets[name]})
	}
	return groups
}

// Keys returns a flat, deduplicated, sorted list of group names present in the
// supplied entries.
func Keys(entries []parser.Entry, separator string) []string {
	if separator == "" {
		separator = "_"
	}
	seen := make(map[string]struct{})
	for _, e := range entries {
		if idx := strings.Index(e.Key, separator); idx > 0 {
			seen[e.Key[:idx]] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
