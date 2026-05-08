// Package trimmer removes unused or stale keys from a .env file by comparing
// it against a reference set of expected keys.
package trimmer

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/parser"
)

// Result holds the outcome of a trim operation.
type Result struct {
	Removed []string
	Kept    []parser.Entry
}

// Trim reads the target .env file, removes any key not present in allowedKeys,
// and writes the pruned file back to disk. It returns a Result describing what
// was removed and what was kept.
func Trim(targetPath string, allowedKeys []string, dryRun bool) (Result, error) {
	entries, err := parser.Parse(targetPath)
	if err != nil {
		return Result{}, fmt.Errorf("trimmer: parse %q: %w", targetPath, err)
	}

	allowed := make(map[string]struct{}, len(allowedKeys))
	for _, k := range allowedKeys {
		allowed[k] = struct{}{}
	}

	var kept []parser.Entry
	var removed []string

	for _, e := range entries {
		if _, ok := allowed[e.Key]; ok {
			kept = append(kept, e)
		} else {
			removed = append(removed, e.Key)
		}
	}

	if !dryRun {
		if err := writeEntries(targetPath, kept); err != nil {
			return Result{}, fmt.Errorf("trimmer: write %q: %w", targetPath, err)
		}
	}

	return Result{Removed: removed, Kept: kept}, nil
}

// writeEntries serialises entries back to a .env file.
func writeEntries(path string, entries []parser.Entry) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, e := range entries {
		line := fmt.Sprintf("%s=%s\n", e.Key, e.Value)
		if _, err := fmt.Fprint(f, line); err != nil {
			return err
		}
	}
	return nil
}
