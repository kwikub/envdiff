// Package interpolator resolves variable references within .env file values.
// It supports ${VAR} and $VAR syntax, resolving from the same env set or
// from a provided override map.
package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

var refPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// FallbackToOS allows falling back to os.Getenv when a key is not found
	// in the provided entries.
	FallbackToOS bool
	// Strict causes Interpolate to return an error when a reference cannot
	// be resolved.
	Strict bool
}

// Interpolate resolves variable references in all entry values.
// Entries are resolved in order; forward references may not resolve on the
// first pass if the referenced key appears later in the slice.
func Interpolate(entries []parser.Entry, overrides map[string]string, opts Options) ([]parser.Entry, error) {
	lookup := buildLookup(entries, overrides)

	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		resolved, err := resolve(e.Value, lookup, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", e.Key, err)
		}
		result[i] = parser.Entry{Key: e.Key, Value: resolved}
	}
	return result, nil
}

func buildLookup(entries []parser.Entry, overrides map[string]string) map[string]string {
	lookup := make(map[string]string, len(entries)+len(overrides))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}
	for k, v := range overrides {
		lookup[k] = v
	}
	return lookup
}

func resolve(value string, lookup map[string]string, opts Options) (string, error) {
	var resolveErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		key := extractKey(match)
		if v, ok := lookup[key]; ok {
			return v
		}
		if opts.FallbackToOS {
			if v, ok := os.LookupEnv(key); ok {
				return v
			}
		}
		if opts.Strict {
			resolveErr = fmt.Errorf("unresolved reference: %s", strings.TrimLeft(match, "${"))
			return match
		}
		return match
	})
	return result, resolveErr
}

func extractKey(match string) string {
	match = strings.TrimPrefix(match, "${") 
	match = strings.TrimSuffix(match, "}")
	match = strings.TrimPrefix(match, "$")
	return match
}
