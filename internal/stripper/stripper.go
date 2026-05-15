// Package stripper removes comments and blank lines from .env files,
// producing a clean, minimal output suitable for deployment or diffing.
package stripper

import (
	"bufio"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls which elements are stripped from the file.
type Options struct {
	RemoveComments   bool
	RemoveBlankLines bool
	RemoveDisabled   bool // lines that are commented-out key=value pairs
}

// DefaultOptions returns a sensible default strip configuration.
func DefaultOptions() Options {
	return Options{
		RemoveComments:   true,
		RemoveBlankLines: true,
		RemoveDisabled:   false,
	}
}

// Strip reads the .env file at path and returns cleaned lines according to opts.
func Strip(path string, opts Options) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			if !opts.RemoveBlankLines {
				out = append(out, line)
			}
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			if opts.RemoveDisabled && isDisabledEntry(trimmed) {
				continue
			}
			if !opts.RemoveComments {
				out = append(out, line)
			}
			continue
		}

		out = append(out, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// StripToEntries is a convenience wrapper that returns parsed entries after
// stripping, delegating actual parsing to the parser package.
func StripToEntries(path string, opts Options) ([]parser.Entry, error) {
	return parser.Parse(path)
}

// isDisabledEntry returns true when a comment line looks like a commented-out
// key=value pair, e.g. "# FOO=bar".
func isDisabledEntry(line string) bool {
	rest := strings.TrimSpace(strings.TrimPrefix(line, "#"))
	return strings.Contains(rest, "=")
}
