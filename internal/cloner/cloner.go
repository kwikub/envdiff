// Package cloner copies entries from one .env file to another,
// optionally filtering by key prefix or allow-list.
package cloner

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls how the clone operation behaves.
type Options struct {
	// Prefix restricts cloning to keys that start with this string.
	Prefix string
	// AllowList restricts cloning to exactly these keys (ignored when empty).
	AllowList []string
	// Overwrite replaces existing keys in the destination file.
	Overwrite bool
	// DryRun skips writing and returns what would have been written.
	DryRun bool
}

// Result holds a summary of the clone operation.
type Result struct {
	Copied  []string
	Skipped []string
}

// Clone reads entries from src and merges matching ones into dst.
func Clone(src, dst string, opts Options) (Result, error) {
	srcEntries, err := parser.Parse(src)
	if err != nil {
		return Result{}, fmt.Errorf("cloner: parse source: %w", err)
	}

	dstEntries, err := parser.Parse(dst)
	if err != nil {
		return Result{}, fmt.Errorf("cloner: parse destination: %w", err)
	}

	allowSet := toSet(opts.AllowList)
	dstMap := toMap(dstEntries)

	var result Result
	for _, e := range srcEntries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if len(allowSet) > 0 && !allowSet[e.Key] {
			continue
		}
		if _, exists := dstMap[e.Key]; exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}
		dstMap[e.Key] = e.Value
		result.Copied = append(result.Copied, e.Key)
	}

	if opts.DryRun {
		return result, nil
	}

	merged := mergeMaps(dstEntries, dstMap)
	if err := writeEntries(dst, merged); err != nil {
		return Result{}, fmt.Errorf("cloner: write destination: %w", err)
	}
	return result, nil
}

func mergeMaps(base []parser.Entry, updated map[string]string) []parser.Entry {
	seen := map[string]bool{}
	out := make([]parser.Entry, 0, len(updated))
	for _, e := range base {
		e.Value = updated[e.Key]
		out = append(out, e)
		seen[e.Key] = true
	}
	for k, v := range updated {
		if !seen[k] {
			out = append(out, parser.Entry{Key: k, Value: v})
		}
	}
	return out
}

func writeEntries(path string, entries []parser.Entry) error {
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
