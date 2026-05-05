package differ

import "strings"

// secretPatterns contains substrings that indicate a key holds a secret value.
var secretPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const maskedValue = "***"

// IsSecret reports whether the given key name looks like it holds a secret.
func IsSecret(key string) bool {
	upper := strings.ToUpper(key)
	for _, pattern := range secretPatterns {
		if strings.Contains(upper, pattern) {
			return true
		}
	}
	return false
}

// MaskSecrets returns a copy of the diff entries with secret values replaced
// by the masked placeholder.
func MaskSecrets(entries []DiffEntry) []DiffEntry {
	out := make([]DiffEntry, len(entries))
	for i, e := range entries {
		if IsSecret(e.Key) {
			if e.BaseVal != "" {
				e.BaseVal = maskedValue
			}
			if e.OtherVal != "" {
				e.OtherVal = maskedValue
			}
		}
		out[i] = e
	}
	return out
}
