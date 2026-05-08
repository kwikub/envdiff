// Package redactor provides functionality to redact sensitive values
// from .env entries before logging, displaying, or exporting them.
package redactor

import (
	"strings"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

const redactedPlaceholder = "[REDACTED]"

// Options controls redaction behaviour.
type Options struct {
	// ExtraKeys are additional key names (case-insensitive) to treat as secret.
	ExtraKeys []string
	// Placeholder overrides the default redaction string.
	Placeholder string
}

// Redact returns a copy of entries with secret values replaced by a
// placeholder. Non-secret entries are returned unchanged.
func Redact(entries []parser.Entry, opts Options) []parser.Entry {
	placeholder := redactedPlaceholder
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}

	extraSet := make(map[string]struct{}, len(opts.ExtraKeys))
	for _, k := range opts.ExtraKeys {
		extraSet[strings.ToUpper(k)] = struct{}{}
	}

	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		if isSecret(e.Key, extraSet) {
			out[i] = parser.Entry{Key: e.Key, Value: placeholder}
		} else {
			out[i] = e
		}
	}
	return out
}

// RedactMap returns a copy of a string map with secret values replaced.
func RedactMap(m map[string]string, opts Options) map[string]string {
	placeholder := redactedPlaceholder
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}

	extraSet := make(map[string]struct{}, len(opts.ExtraKeys))
	for _, k := range opts.ExtraKeys {
		extraSet[strings.ToUpper(k)] = struct{}{}
	}

	out := make(map[string]string, len(m))
	for k, v := range m {
		if isSecret(k, extraSet) {
			out[k] = placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

func isSecret(key string, extra map[string]struct{}) bool {
	if differ.IsSecret(key) {
		return true
	}
	_, ok := extra[strings.ToUpper(key)]
	return ok
}
