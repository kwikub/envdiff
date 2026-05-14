// Package highlighter provides syntax-aware colorised output for .env file diffs.
package highlighter

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// ANSI colour codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// Options controls highlighter behaviour.
type Options struct {
	// NoColor disables ANSI escape sequences (e.g. when writing to a file).
	NoColor bool
	// ShowUnchanged includes UNCHANGED entries in the output.
	ShowUnchanged bool
}

// Highlight renders a slice of differ.Result lines as a colourised diff string.
func Highlight(results []differ.Result, opts Options) string {
	var sb strings.Builder
	for _, r := range results {
		line := formatResult(r, opts)
		if line == "" {
			continue
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func formatResult(r differ.Result, opts Options) string {
	switch r.Status {
	case differ.Added:
		return colorise(opts, colorGreen, fmt.Sprintf("+ %s=%s", r.Key, r.NewValue))
	case differ.Removed:
		return colorise(opts, colorRed, fmt.Sprintf("- %s=%s", r.Key, r.OldValue))
	case differ.Modified:
		oldLine := colorise(opts, colorRed, fmt.Sprintf("- %s=%s", r.Key, r.OldValue))
		newLine := colorise(opts, colorGreen, fmt.Sprintf("+ %s=%s", r.Key, r.NewValue))
		return oldLine + "\n" + newLine
	case differ.Unchanged:
		if !opts.ShowUnchanged {
			return ""
		}
		return colorise(opts, colorGray, fmt.Sprintf("  %s=%s", r.Key, r.OldValue))
	default:
		return ""
	}
}

func colorise(opts Options, code, text string) string {
	if opts.NoColor {
		return text
	}
	return code + text + colorReset
}

// Summary returns a one-line colourised summary of diff counts.
func Summary(results []differ.Result, opts Options) string {
	var added, removed, modified int
	for _, r := range results {
		switch r.Status {
		case differ.Added:
			added++
		case differ.Removed:
			removed++
		case differ.Modified:
			modified++
		}
	}
	parts := []string{
		colorise(opts, colorGreen, fmt.Sprintf("+%d added", added)),
		colorise(opts, colorRed, fmt.Sprintf("-%d removed", removed)),
		colorise(opts, colorYellow, fmt.Sprintf("~%d modified", modified)),
	}
	return strings.Join(parts, "  ")
}
