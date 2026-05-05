// Package reconciler provides functionality to compare two .env files and
// produce actionable suggestions to bring a target environment file in sync
// with a source environment file.
//
// It builds on top of the differ and parser packages:
//   - parser.Parse reads and tokenises .env files into key/value entries.
//   - differ.Diff computes the delta between two sets of entries.
//   - reconciler.Reconcile converts that delta into typed Suggestion values
//     (add, remove, update, keep) that a CLI or other consumer can act on.
//
// Secret masking is applied automatically via differ.IsSecret when
// Format is called with maskSecrets=true, ensuring sensitive values are
// never printed in plain text.
package reconciler
