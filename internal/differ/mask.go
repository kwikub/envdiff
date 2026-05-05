package differ

import (
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// secretPatterns holds substrings that indicate a key likely holds a secret value.
var secretPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"ACCESS_KEY",
	"AUTH",
	"CREDENTIAL",
}

// IsSecret reports whether the given environment variable key is likely
// to contain a sensitive / secret value based on common naming patterns.
func IsSecret(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range secretPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// MaskSecrets returns a copy of the provided entries with secret values
// replaced by "***". The original slice is not modified.
func MaskSecrets(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		if IsSecret(e.Key) {
			e.Value = "***"
		}
		result[i] = e
	}
	return result
}
