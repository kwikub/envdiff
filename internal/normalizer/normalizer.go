// Package normalizer provides utilities for normalising .env entry keys and
// values into a canonical form — trimming whitespace, standardising key
// casing, and stripping redundant quotes.
package normalizer

import (
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls how normalisation is applied.
type Options struct {
	// UppercaseKeys converts all key names to UPPER_CASE.
	UppercaseKeys bool
	// TrimValues strips leading and trailing whitespace from values.
	TrimValues bool
	// StripQuotes removes surrounding single or double quotes from values.
	StripQuotes bool
}

// Normalize applies the given options to each entry and returns a new slice
// with the normalised entries. The original slice is never mutated.
func Normalize(entries []parser.Entry, opts Options) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = normaliseEntry(e, opts)
	}
	return out
}

func normaliseEntry(e parser.Entry, opts Options) parser.Entry {
	key := e.Key
	val := e.Value

	if opts.UppercaseKeys {
		key = strings.ToUpper(key)
	}

	if opts.TrimValues {
		val = strings.TrimSpace(val)
	}

	if opts.StripQuotes {
		val = stripQuotes(val)
	}

	return parser.Entry{Key: key, Value: val}
}

// stripQuotes removes a matching pair of surrounding quotes (single or double)
// from s, if present. It does not strip mismatched or nested quotes.
func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}
