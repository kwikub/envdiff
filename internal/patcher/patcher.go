// Package patcher applies a set of key-value patches to an existing .env file,
// updating, adding, or removing entries as directed.
package patcher

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Op represents the type of patch operation.
type Op string

const (
	OpSet    Op = "set"    // add or update a key
	OpDelete Op = "delete" // remove a key
)

// Patch describes a single change to apply.
type Patch struct {
	Op    Op
	Key   string
	Value string
}

// Result holds the patched entries and a summary of changes applied.
type Result struct {
	Entries []parser.Entry
	Applied []Patch
	Skipped []Patch // ops that had no effect (e.g. delete of missing key)
}

// Apply applies the given patches to entries parsed from file at path.
// It returns a Result containing the modified entries and change summary.
func Apply(path string, patches []Patch) (*Result, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("patcher: parse %q: %w", path, err)
	}

	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	var applied, skipped []Patch

	for _, p := range patches {
		switch p.Op {
		case OpSet:
			if idx, ok := index[p.Key]; ok {
				entries[idx].Value = p.Value
			} else {
				newEntry := parser.Entry{Key: p.Key, Value: p.Value}
				index[p.Key] = len(entries)
				entries = append(entries, newEntry)
			}
			applied = append(applied, p)
		case OpDelete:
			if _, ok := index[p.Key]; ok {
				entries = removeKey(entries, p.Key)
				// rebuild index after removal
				index = make(map[string]int, len(entries))
				for i, e := range entries {
					index[e.Key] = i
				}
				applied = append(applied, p)
			} else {
				skipped = append(skipped, p)
			}
		default:
			return nil, fmt.Errorf("patcher: unknown op %q for key %q", p.Op, p.Key)
		}
	}

	return &Result{Entries: entries, Applied: applied, Skipped: skipped}, nil
}

// Format serialises patched entries back to .env file content.
func Format(entries []parser.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return sb.String()
}

func removeKey(entries []parser.Entry, key string) []parser.Entry {
	out := entries[:0:len(entries)]
	for _, e := range entries {
		if e.Key != key {
			out = append(out, e)
		}
	}
	return out
}
