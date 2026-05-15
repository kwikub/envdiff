// Package cloner provides functionality for copying environment variable
// entries from one .env file into another.
//
// Entries can be filtered by key prefix or an explicit allow-list.
// Existing keys in the destination are preserved by default; pass
// Options.Overwrite to replace them. Use Options.DryRun to preview
// what would be copied without touching the filesystem.
//
// Example:
//
//	res, err := cloner.Clone(".env.production", ".env.staging", cloner.Options{
//		Prefix:    "APP_",
//		Overwrite: false,
//		DryRun:    false,
//	})
package cloner
