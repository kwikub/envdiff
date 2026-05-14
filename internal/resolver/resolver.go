// Package resolver provides functionality for resolving the final effective
// value of a key across multiple .env files, applying precedence rules.
package resolver

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Strategy controls how values are resolved when a key appears in multiple files.
type Strategy string

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = "first"
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast Strategy = "last"
)

// Result holds the resolved value for a single key along with provenance info.
type Result struct {
	Key      string
	Value    string
	Source   string // file path the winning value came from
	Overridden bool // true if the key appeared in more than one file
}

// Resolve reads the given env files in order and returns one Result per unique
// key, honouring the chosen precedence strategy.
func Resolve(files []string, strategy Strategy) ([]Result, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("resolver: at least one file is required")
	}
	if strategy != StrategyFirst && strategy != StrategyLast {
		return nil, fmt.Errorf("resolver: unknown strategy %q", strategy)
	}

	type record struct {
		value  string
		source string
		count  int
	}

	// Preserve insertion order via a slice of keys.
	order := []string{}
	seen := map[string]*record{}

	for _, path := range files {
		entries, err := parser.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("resolver: parsing %s: %w", path, err)
		}
		for _, e := range entries {
			rec, exists := seen[e.Key]
			if !exists {
				order = append(order, e.Key)
				seen[e.Key] = &record{value: e.Value, source: path, count: 1}
				continue
			}
			rec.count++
			if strategy == StrategyLast {
				rec.value = e.Value
				rec.source = path
			}
		}
	}

	results := make([]Result, 0, len(order))
	for _, k := range order {
		rec := seen[k]
		results = append(results, Result{
			Key:        k,
			Value:      rec.value,
			Source:     rec.source,
			Overridden: rec.count > 1,
		})
	}
	return results, nil
}
