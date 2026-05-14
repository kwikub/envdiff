// Package differ provides utilities for computing differences between two sets
// of .env entries.
//
// The Diff function returns a list of DiffResult values describing keys that
// were added, removed, modified, or unchanged between a base and target file.
//
// The Compare function performs a side-by-side comparison of two entry slices,
// returning a CompareResult for every key found in either side, annotated with
// a status of "match", "mismatch", "left_only", or "right_only".
//
// Secret masking helpers are provided via MaskSecrets and IsSecret to prevent
// sensitive values from appearing in diff output.
package differ
