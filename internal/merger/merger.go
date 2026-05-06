// Package merger provides functionality to merge multiple .env files,
// with later files taking precedence over earlier ones.
package merger

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// MergeStrategy controls how conflicts between files are resolved.
type MergeStrategy int

const (
	// StrategyLast means the last file's value wins on conflict.
	StrategyLast MergeStrategy = iota
	// StrategyFirst means the first file's value wins on conflict.
	StrategyFirst
)

// Result holds the merged entries and metadata about the merge.
type Result struct {
	Entries  []parser.Entry
	Sources  map[string]string // key -> source filename
	Conflicts []Conflict
}

// Conflict records a key that appeared in more than one file.
type Conflict struct {
	Key    string
	Values []string // values in order of files
	Files  []string // filenames in order
}

// Merge combines entries from multiple .env files.
// Files are processed in order; strategy determines which value wins on conflict.
func Merge(files []string, strategy MergeStrategy) (*Result, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("merger: no files provided")
	}

	result := &Result{
		Sources:   make(map[string]string),
		Conflicts: []Conflict{},
	}

	// Track order of keys and conflicts per key.
	conflictMap := make(map[string]*Conflict)
	merged := make(map[string]parser.Entry)
	var keyOrder []string

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("merger: open %q: %w", file, err)
		}
		entries, err := parser.Parse(f)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("merger: parse %q: %w", file, err)
		}

		for _, e := range entries {
			if existing, exists := merged[e.Key]; exists {
				// Record conflict.
				if c, ok := conflictMap[e.Key]; ok {
					c.Values = append(c.Values, e.Value)
					c.Files = append(c.Files, file)
				} else {
					conflictMap[e.Key] = &Conflict{
						Key:    e.Key,
						Values: []string{existing.Value, e.Value},
						Files:  []string{result.Sources[e.Key], file},
					}
				}
				if strategy == StrategyLast {
					merged[e.Key] = e
					result.Sources[e.Key] = file
				}
			} else {
				merged[e.Key] = e
				result.Sources[e.Key] = file
				keyOrder = append(keyOrder, e.Key)
			}
		}
	}

	for _, key := range keyOrder {
		result.Entries = append(result.Entries, merged[key])
	}
	for _, c := range conflictMap {
		result.Conflicts = append(result.Conflicts, *c)
	}

	_ = strings.TrimSpace // imported for potential future use
	return result, nil
}
