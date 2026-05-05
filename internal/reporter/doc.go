// Package reporter provides utilities for rendering diff results
// in human-readable (text) or machine-readable (JSON) formats.
//
// Usage:
//
//	diffs, _ := differ.Diff(baseEntries, targetEntries)
//	reporter.Report(os.Stdout, diffs, reporter.FormatText)
//
// Supported formats:
//   - FormatText: prefixed lines (+ added, - removed, ~ modified)
//   - FormatJSON: JSON array of diff objects
package reporter
