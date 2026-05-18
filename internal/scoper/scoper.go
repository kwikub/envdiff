// Package scoper restricts env entries to a named scope (e.g. "production",
// "staging") by reading or writing a scope tag embedded in key comments.
package scoper

import (
	"fmt"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// Entry mirrors parser.Entry for convenience.
type Entry = parser.Entry

// Options controls how scoping is applied.
type Options struct {
	// Scope is the target scope name (e.g. "production").
	Scope string
	// TagPrefix is the inline comment prefix used to mark scope.
	// Defaults to "@scope:".
	TagPrefix string
	// Strict causes Scope to return an error when an entry carries a
	// different scope tag rather than silently dropping it.
	Strict bool
}

func tagPrefix(o Options) string {
	if o.TagPrefix != "" {
		return o.TagPrefix
	}
	return "@scope:"
}

// Scope filters entries to those that either carry no scope tag or whose
// scope tag matches opts.Scope. When opts.Strict is true, entries tagged
// with a different scope return an error.
func Scope(entries []Entry, opts Options) ([]Entry, error) {
	if opts.Scope == "" {
		return entries, nil
	}
	prefix := tagPrefix(opts)
	var out []Entry
	for _, e := range entries {
		scope, tagged := extractScope(e.Comment, prefix)
		if !tagged {
			out = append(out, e)
			continue
		}
		if scope == opts.Scope {
			out = append(out, e)
			continue
		}
		if opts.Strict {
			return nil, fmt.Errorf("scoper: key %q is tagged for scope %q, not %q", e.Key, scope, opts.Scope)
		}
	}
	return out, nil
}

// Tag annotates entries whose keys are in keys with a scope comment tag.
// Existing scope tags on those entries are replaced.
func Tag(entries []Entry, scope string, keys []string, tagPrefix string) []Entry {
	if tagPrefix == "" {
		tagPrefix = "@scope:"
	}
	keySet := toSet(keys)
	out := make([]Entry, len(entries))
	for i, e := range entries {
		if keySet[e.Key] {
			e.Comment = setScope(e.Comment, scope, tagPrefix)
		}
		out[i] = e
	}
	return out
}

// extractScope parses a comment string and returns the scope value and
// whether a scope tag was found.
func extractScope(comment, prefix string) (string, bool) {
	idx := strings.Index(comment, prefix)
	if idx == -1 {
		return "", false
	}
	rest := comment[idx+len(prefix):]
	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return "", false
	}
	return fields[0], true
}

// setScope replaces or appends a scope tag in a comment string.
func setScope(comment, scope, prefix string) string {
	if idx := strings.Index(comment, prefix); idx != -1 {
		before := comment[:idx]
		rest := comment[idx+len(prefix):]
		fields := strings.Fields(rest)
		if len(fields) > 0 {
			fields[0] = scope
		} else {
			fields = []string{scope}
		}
		return strings.TrimRight(before, " ") + " " + prefix + strings.Join(fields, " ")
	}
	if comment == "" {
		return prefix + scope
	}
	return comment + " " + prefix + scope
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
