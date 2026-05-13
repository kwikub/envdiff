// Package filter provides utilities for filtering .env entries
// by key patterns, prefixes, or custom predicates.
package filter

import (
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls how entries are filtered.
type Options struct {
	// Prefix keeps only entries whose keys start with the given prefix.
	Prefix string
	// Keys keeps only entries whose keys are in the provided set.
	// If empty, this constraint is not applied.
	Keys []string
	// Exclude removes entries whose keys are in the provided set.
	Exclude []string
	// SecretsOnly keeps only entries that are considered secrets.
	SecretsOnly bool
}

// Filter returns a subset of entries that match the given options.
// All non-empty constraints are ANDed together.
func Filter(entries []parser.Entry, opts Options) []parser.Entry {
	allowSet := toSet(opts.Keys)
	excludeSet := toSet(opts.Exclude)

	var result []parser.Entry
	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if len(allowSet) > 0 && !allowSet[e.Key] {
			continue
		}
		if excludeSet[e.Key] {
			continue
		}
		if opts.SecretsOnly && !isSecret(e.Key) {
			continue
		}
		result = append(result, e)
	}
	return result
}

// isSecret returns true if the key name suggests it holds a secret value.
func isSecret(key string) bool {
	upper := strings.ToUpper(key)
	secretTerms := []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL"}
	for _, term := range secretTerms {
		if strings.Contains(upper, term) {
			return true
		}
	}
	return false
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
