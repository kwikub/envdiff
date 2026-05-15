// Package differ provides utilities for comparing two sets of .env entries.
//
// Diff returns a list of DiffEntry values describing keys that were added,
// removed, modified, or unchanged between a base and a target environment.
//
// Compare performs a symmetric key-by-key comparison and reports the status
// of every key seen in either input, making it suitable for side-by-side
// inspection of two arbitrary env files.
//
// MaskSecrets and FilterSecrets help redact sensitive values before the
// diff results are displayed or written to logs.
package differ
