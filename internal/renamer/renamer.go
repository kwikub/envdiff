// Package renamer provides utilities for renaming keys across .env files.
package renamer

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/parser"
)

// RenameResult describes the outcome of a rename operation on a single file.
type RenameResult struct {
	File    string
	OldKey  string
	NewKey  string
	Renamed bool
	Err     error
}

// Rename renames oldKey to newKey in each of the provided .env files.
// It returns one RenameResult per file. Files where the key is not found
// are recorded with Renamed=false but no error.
func Rename(files []string, oldKey, newKey string) []RenameResult {
	results := make([]RenameResult, 0, len(files))
	for _, f := range files {
		res := renameInFile(f, oldKey, newKey)
		results = append(results, res)
	}
	return results
}

func renameInFile(path, oldKey, newKey string) RenameResult {
	result := RenameResult{File: path, OldKey: oldKey, NewKey: newKey}

	entries, err := parser.Parse(path)
	if err != nil {
		result.Err = fmt.Errorf("parse %s: %w", path, err)
		return result
	}

	renamed := false
	for i, e := range entries {
		if e.Key == oldKey {
			entries[i].Key = newKey
			renamed = true
		}
	}

	if !renamed {
		return result
	}

	if err := writeEntries(path, entries); err != nil {
		result.Err = fmt.Errorf("write %s: %w", path, err)
		return result
	}

	result.Renamed = true
	return result
}

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
