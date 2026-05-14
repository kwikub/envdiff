// Package promoter copies approved env entries from one environment file to another,
// optionally masking secrets and skipping keys already present in the target.
package promoter

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/parser"
)

// Options controls promotion behaviour.
type Options struct {
	// SkipExisting prevents overwriting keys already defined in the target.
	SkipExisting bool
	// Keys is an optional allow-list; when non-empty only these keys are promoted.
	Keys []string
	// DryRun reports what would change without writing the target file.
	DryRun bool
}

// Result describes the outcome of a single key promotion.
type Result struct {
	Key      string
	Action   string // "promoted", "skipped", "dry-run"
	OldValue string // empty when Action == "promoted" and key was absent
	NewValue string
}

// Promote reads entries from src, applies them to dst according to opts,
// writes the updated dst file (unless DryRun), and returns a result per key.
func Promote(src, dst string, opts Options) ([]Result, error) {
	srcEntries, err := parser.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("promoter: parse src %q: %w", src, err)
	}
	dstEntries, err := parser.Parse(dst)
	if err != nil {
		return nil, fmt.Errorf("promoter: parse dst %q: %w", dst, err)
	}

	allowSet := toSet(opts.Keys)
	dstMap := toMap(dstEntries)

	var results []Result
	updated := make([]parser.Entry, len(dstEntries))
	copy(updated, dstEntries)

	for _, e := range srcEntries {
		if len(allowSet) > 0 && !allowSet[e.Key] {
			continue
		}
		old, exists := dstMap[e.Key]
		if exists && opts.SkipExisting {
			results = append(results, Result{Key: e.Key, Action: "skipped", OldValue: old, NewValue: old})
			continue
		}
		if opts.DryRun {
			results = append(results, Result{Key: e.Key, Action: "dry-run", OldValue: old, NewValue: e.Value})
			continue
		}
		updated = upsert(updated, e)
		results = append(results, Result{Key: e.Key, Action: "promoted", OldValue: old, NewValue: e.Value})
	}

	if !opts.DryRun {
		if err := writeEntries(dst, updated); err != nil {
			return nil, fmt.Errorf("promoter: write dst %q: %w", dst, err)
		}
	}
	return results, nil
}

func upsert(entries []parser.Entry, e parser.Entry) []parser.Entry {
	for i, ex := range entries {
		if ex.Key == e.Key {
			entries[i].Value = e.Value
			return entries
		}
	}
	return append(entries, e)
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func writeEntries(path string, entries []parser.Entry) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, e := range entries {
		_, err := fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
