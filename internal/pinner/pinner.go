// Package pinner provides functionality to pin (lock) the current values
// of .env entries into a lockfile, similar to a dependency lockfile.
package pinner

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"time"

	"github.com/user/envdiff/internal/parser"
)

// PinFile represents a lockfile containing pinned env key-value pairs.
type PinFile struct {
	PinnedAt time.Time         `json:"pinned_at"`
	Source   string            `json:"source"`
	Entries  map[string]string `json:"entries"`
}

// Pin reads the given .env file and writes a lockfile to dest.
func Pin(src, dest string) (*PinFile, error) {
	entries, err := parser.Parse(src)
	if err != nil {
		return nil, err
	}

	pf := &PinFile{
		PinnedAt: time.Now().UTC(),
		Source:   src,
		Entries:  make(map[string]string, len(entries)),
	}
	for _, e := range entries {
		pf.Entries[e.Key] = e.Value
	}

	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return nil, err
	}
	return pf, os.WriteFile(dest, data, 0o600)
}

// Load reads a previously written lockfile from path.
func Load(path string) (*PinFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, err
	}
	return &pf, nil
}

// Verify compares the current .env file against the lockfile.
// It returns a list of keys whose values have drifted from the pinned state.
func Verify(src, lockPath string) ([]string, error) {
	pf, err := Load(lockPath)
	if err != nil {
		return nil, err
	}

	entries, err := parser.Parse(src)
	if err != nil {
		return nil, err
	}

	current := make(map[string]string, len(entries))
	for _, e := range entries {
		current[e.Key] = e.Value
	}

	if len(pf.Entries) == 0 {
		return nil, errors.New("pinner: lockfile contains no entries")
	}

	var drifted []string
	for key, pinnedVal := range pf.Entries {
		if curVal, ok := current[key]; !ok || curVal != pinnedVal {
			drifted = append(drifted, key)
		}
	}
	sort.Strings(drifted)
	return drifted, nil
}
