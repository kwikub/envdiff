// Package profiler provides environment profile management, allowing
// named profiles (e.g. "staging", "production") to be stored and retrieved.
package profiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// Profile represents a named collection of environment entries.
type Profile struct {
	Name    string           `json:"name"`
	Entries []parser.Entry   `json:"entries"`
}

// Save writes a profile to a JSON file in the given directory.
func Save(dir string, profile Profile) error {
	if profile.Name == "" {
		return fmt.Errorf("profile name must not be empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating profile directory: %w", err)
	}
	path := filepath.Join(dir, profile.Name+".json")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating profile file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(profile); err != nil {
		return fmt.Errorf("encoding profile: %w", err)
	}
	return nil
}

// Load reads a named profile from the given directory.
func Load(dir, name string) (Profile, error) {
	path := filepath.Join(dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return Profile{}, fmt.Errorf("reading profile %q: %w", name, err)
	}
	var p Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return Profile{}, fmt.Errorf("parsing profile %q: %w", name, err)
	}
	return p, nil
}

// List returns the names of all profiles stored in the given directory.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading profile directory: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
