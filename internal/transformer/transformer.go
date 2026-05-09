// Package transformer applies key/value transformations to env entries.
// Supported transforms: uppercase keys, lowercase values, prefix addition,
// prefix stripping, and trimming whitespace from values.
package transformer

import (
	"fmt"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// Transform defines a named transformation to apply.
type Transform string

const (
	UppercaseKeys  Transform = "uppercase-keys"
	LowercaseKeys  Transform = "lowercase-keys"
	TrimValues     Transform = "trim-values"
	AddPrefix      Transform = "add-prefix"
	StripPrefix    Transform = "strip-prefix"
)

// Options configures the transformation run.
type Options struct {
	Transforms []Transform
	Prefix     string // used by add-prefix and strip-prefix
}

// Apply runs all configured transforms against entries and returns new entries.
// Original entries are never mutated.
func Apply(entries []parser.Entry, opts Options) ([]parser.Entry, error) {
	result := make([]parser.Entry, len(entries))
	copy(result, entries)

	for _, t := range opts.Transforms {
		var err error
		result, err = applyOne(result, t, opts)
		if err != nil {
			return nil, fmt.Errorf("transform %q: %w", t, err)
		}
	}
	return result, nil
}

func applyOne(entries []parser.Entry, t Transform, opts Options) ([]parser.Entry, error) {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		switch t {
		case UppercaseKeys:
			e.Key = strings.ToUpper(e.Key)
		case LowercaseKeys:
			e.Key = strings.ToLower(e.Key)
		case TrimValues:
			e.Value = strings.TrimSpace(e.Value)
		case AddPrefix:
			if opts.Prefix == "" {
				return nil, fmt.Errorf("prefix must not be empty")
			}
			if !strings.HasPrefix(e.Key, opts.Prefix) {
				e.Key = opts.Prefix + e.Key
			}
		case StripPrefix:
			if opts.Prefix == "" {
				return nil, fmt.Errorf("prefix must not be empty")
			}
			e.Key = strings.TrimPrefix(e.Key, opts.Prefix)
		default:
			return nil, fmt.Errorf("unknown transform %q", t)
		}
		out[i] = e
	}
	return out, nil
}
