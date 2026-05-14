// Package sanitizer provides utilities for sanitizing .env file values
// by normalising whitespace, removing unsafe characters, and enforcing
// key naming conventions before writing or exporting.
package sanitizer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Options controls which sanitization passes are applied.
type Options struct {
	TrimValues     bool // strip leading/trailing whitespace from values
	UppercaseKeys  bool // convert keys to UPPER_CASE
	StripQuotes    bool // remove surrounding quotes from values
	RejectInvalid  bool // return an error for keys that violate naming rules
}

// Result holds a sanitized entry alongside any warning produced.
type Result struct {
	Entry   parser.Entry
	Warning string
}

// Sanitize applies the requested passes to each entry and returns the
// transformed results. If RejectInvalid is set and a key is malformed,
// an error is returned immediately.
func Sanitize(entries []parser.Entry, opts Options) ([]Result, error) {
	out := make([]Result, 0, len(entries))

	for _, e := range entries {
		var warn string

		key := e.Key
		val := e.Value

		if opts.UppercaseKeys {
			key = strings.ToUpper(key)
		}

		if opts.RejectInvalid && !validKeyRe.MatchString(key) {
			return nil, fmt.Errorf("sanitizer: invalid key name %q", key)
		}

		if !validKeyRe.MatchString(key) {
			warn = fmt.Sprintf("key %q does not match recommended naming convention", key)
		}

		if opts.StripQuotes {
			val = stripSurroundingQuotes(val)
		}

		if opts.TrimValues {
			val = strings.TrimSpace(val)
		}

		out = append(out, Result{
			Entry:   parser.Entry{Key: key, Value: val},
			Warning: warn,
		})
	}

	return out, nil
}

func stripSurroundingQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
