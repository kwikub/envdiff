// Package splitter splits a flat list of env entries into multiple files
// based on a grouping strategy (prefix, tag, or custom mapping).
package splitter

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Strategy controls how entries are assigned to output files.
type Strategy string

const (
	StrategyPrefix Strategy = "prefix" // group by KEY prefix before "_"
	StrategyAlpha  Strategy = "alpha"  // split alphabetically into N buckets
)

// Options configures a Split operation.
type Options struct {
	Strategy  Strategy
	// Mapping overrides automatic grouping: map[outputFile][]keyPrefix
	Mapping   map[string][]string
	// Buckets is used with StrategyAlpha to control number of output files.
	Buckets   int
}

// Split reads src, partitions its entries according to opts, and writes each
// partition to the corresponding path in destDir.  It returns a map of
// output-file-name -> entry count.
func Split(src string, destDir string, opts Options) (map[string]int, error) {
	entries, err := parser.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("splitter: parse %s: %w", src, err)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, fmt.Errorf("splitter: mkdir %s: %w", destDir, err)
	}

	var groups map[string][]parser.Entry

	switch {
	case len(opts.Mapping) > 0:
		groups = applyMapping(entries, opts.Mapping)
	case opts.Strategy == StrategyAlpha:
		buckets := opts.Buckets
		if buckets < 2 {
			buckets = 2
		}
		groups = splitAlpha(entries, buckets)
	default:
		groups = splitByPrefix(entries)
	}

	counts := make(map[string]int, len(groups))
	for name, grp := range groups {
		path := destDir + "/" + name
		if err := writeGroup(path, grp); err != nil {
			return nil, err
		}
		counts[name] = len(grp)
	}
	return counts, nil
}

func splitByPrefix(entries []parser.Entry) map[string][]parser.Entry {
	out := map[string][]parser.Entry{}
	for _, e := range entries {
		prefix := groupOf(e.Key)
		fileName := strings.ToLower(prefix) + ".env"
		out[fileName] = append(out[fileName], e)
	}
	return out
}

func groupOf(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}

func splitAlpha(entries []parser.Entry, buckets int) map[string][]parser.Entry {
	out := map[string][]parser.Entry{}
	for _, e := range entries {
		if len(e.Key) == 0 {
			continue
		}
		idx := int(e.Key[0]) % buckets
		name := fmt.Sprintf("part%02d.env", idx)
		out[name] = append(out[name], e)
	}
	return out
}

func applyMapping(entries []parser.Entry, mapping map[string][]string) map[string][]parser.Entry {
	out := map[string][]parser.Entry{}
	for _, e := range entries {
		assigned := false
		for file, prefixes := range mapping {
			for _, p := range prefixes {
				if strings.HasPrefix(e.Key, p) {
					out[file] = append(out[file], e)
					assigned = true
					break
				}
			}
			if assigned {
				break
			}
		}
		if !assigned {
			out["other.env"] = append(out["other.env"], e)
		}
	}
	return out
}

func writeGroup(path string, entries []parser.Entry) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("splitter: create %s: %w", path, err)
	}
	defer f.Close()
	for _, e := range entries {
		if _, err := fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value); err != nil {
			return err
		}
	}
	return nil
}
