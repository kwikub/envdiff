// Package template provides functionality to generate .env template files
// from existing env entries, replacing values with placeholder descriptions.
package template

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

// Options controls template generation behaviour.
type Options struct {
	// MaskSecrets replaces secret values with a redacted placeholder.
	MaskSecrets bool
	// AddComments adds a descriptive comment above each key.
	AddComments bool
}

// Generate reads a .env file and writes a template version to outPath.
// Values are replaced with empty strings or typed placeholders.
func Generate(srcPath, outPath string, opts Options) error {
	entries, err := parser.Parse(srcPath)
	if err != nil {
		return fmt.Errorf("template: parse %q: %w", srcPath, err)
	}

	var sb strings.Builder
	for _, e := range entries {
		if opts.AddComments {
			sb.WriteString(fmt.Sprintf("# %s\n", describeKey(e.Key)))
		}
		value := ""
		if opts.MaskSecrets && differ.IsSecret(e.Key) {
			value = "<secret>"
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, value))
	}

	if err := os.WriteFile(outPath, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("template: write %q: %w", outPath, err)
	}
	return nil
}

// describeKey returns a human-readable hint for a given key name.
func describeKey(key string) string {
	lower := strings.ToLower(key)
	switch {
	case strings.Contains(lower, "url") || strings.Contains(lower, "uri"):
		return fmt.Sprintf("%s — connection URL", key)
	case strings.Contains(lower, "port"):
		return fmt.Sprintf("%s — port number", key)
	case strings.Contains(lower, "host"):
		return fmt.Sprintf("%s — hostname or IP", key)
	case differ.IsSecret(key):
		return fmt.Sprintf("%s — secret value (required)", key)
	default:
		return fmt.Sprintf("%s — set your value here", key)
	}
}
