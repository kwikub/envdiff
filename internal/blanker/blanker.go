// Package blanker replaces the values of matching .env entries with empty
// strings, leaving keys intact. This is useful for generating sanitised
// skeleton files that can be committed to source control.
package blanker

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls which entries are blanked.
type Options struct {
	// Keys is an explicit allow-list of key names to blank. When non-empty,
	// only these keys are affected.
	Keys []string
	// SecretsOnly blanks only keys that look like secrets (password, token, etc.).
	SecretsOnly bool
	// Placeholder is written as the value instead of an empty string.
	// Defaults to "" when left unset.
	Placeholder string
}

// Blank reads the .env file at path, blanks matching values according to opts,
// and returns the modified entries.
func Blank(path string, opts Options) ([]parser.Entry, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("blanker: parse %q: %w", path, err)
	}

	allowSet := toSet(opts.Keys)

	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		if shouldBlank(e.Key, allowSet, opts.SecretsOnly) {
			e.Value = opts.Placeholder
		}
		result[i] = e
	}
	return result, nil
}

// Format serialises entries back to .env file content.
func Format(entries []parser.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		if e.Value == "" {
			fmt.Fprintf(&sb, "%s=\n", e.Key)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
		}
	}
	return sb.String()
}

func shouldBlank(key string, allowSet map[string]struct{}, secretsOnly bool) bool {
	if len(allowSet) > 0 {
		_, ok := allowSet[key]
		return ok
	}
	if secretsOnly {
		return isSecret(key)
	}
	return true
}

func isSecret(key string) bool {
	lower := strings.ToLower(key)
	for _, kw := range []string{"secret", "password", "passwd", "token", "apikey", "api_key", "private", "credential"} {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func toSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}
