// Package rotator provides utilities for rotating secret values in .env files,
// generating new values for secret keys while preserving non-secret entries.
package rotator

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Options controls rotation behaviour.
type Options struct {
	// Keys is an explicit list of keys to rotate. If empty, all secret keys are rotated.
	Keys []string
	// Length is the byte-length of the generated secret (hex-encoded output is 2x this).
	// Defaults to 16 (32-char hex string).
	Length int
	// DryRun reports what would change without writing to disk.
	DryRun bool
}

// Result describes a single rotated key.
type Result struct {
	Key      string
	OldValue string
	NewValue string
}

// Rotate reads the .env file at path, replaces secret values (or the
// explicitly listed keys) with freshly generated random hex strings, writes
// the result back to disk (unless DryRun is set), and returns the list of
// rotations performed.
func Rotate(path string, opts Options) ([]Result, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("rotator: parse %q: %w", path, err)
	}

	allowSet := toSet(opts.Keys)
	length := opts.Length
	if length <= 0 {
		length = 16
	}

	var results []Result
	for i, e := range entries {
		if !shouldRotate(e.Key, allowSet) {
			continue
		}
		newVal, err := generateHex(length)
		if err != nil {
			return nil, fmt.Errorf("rotator: generate value for %q: %w", e.Key, err)
		}
		results = append(results, Result{Key: e.Key, OldValue: e.Value, NewValue: newVal})
		entries[i].Value = newVal
	}

	if !opts.DryRun {
		if err := writeEntries(path, entries); err != nil {
			return nil, fmt.Errorf("rotator: write %q: %w", path, err)
		}
	}

	return results, nil
}

func shouldRotate(key string, allowSet map[string]struct{}) bool {
	if len(allowSet) > 0 {
		_, ok := allowSet[key]
		return ok
	}
	return isSecret(key)
}

func isSecret(key string) bool {
	secretSuffixes := []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS", "PRIVATE", "CREDENTIAL"}
	for _, s := range secretSuffixes {
		if len(key) >= len(s) && key[len(key)-len(s):] == s {
			return true
		}
	}
	return false
}

func generateHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func toSet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}
