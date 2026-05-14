// Package promoter provides functionality to promote env entries from one
// environment file (e.g. staging) to another (e.g. production).
//
// Promotion respects an optional allow-list of keys, can skip keys that
// already exist in the target, and supports a dry-run mode that reports
// changes without writing to disk.
//
// Example usage:
//
//	results, err := promoter.Promote(".env.staging", ".env.production",
//		promoter.Options{SkipExisting: true, Keys: []string{"API_URL"}})
package promoter
