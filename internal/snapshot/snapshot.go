// Package snapshot provides functionality to capture and compare
// .env file snapshots over time, enabling drift detection.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/parser"
)

// Snapshot represents a point-in-time capture of an .env file.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Entries   []parser.Entry    `json:"entries"`
}

// Take reads the given .env file and returns a Snapshot.
func Take(path string) (*Snapshot, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: parse %q: %w", path, err)
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    path,
		Entries:   entries,
	}, nil
}

// Save writes the snapshot as JSON to the given output path.
func Save(snap *Snapshot, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("snapshot: create %q: %w", dest, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a previously saved snapshot from a JSON file.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %q: %w", path, err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &snap, nil
}

// ToMap converts snapshot entries to a key→value map.
func (s *Snapshot) ToMap() map[string]string {
	m := make(map[string]string, len(s.Entries))
	for _, e := range s.Entries {
		m[e.Key] = e.Value
	}
	return m
}
