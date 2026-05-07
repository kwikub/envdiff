// Package renamer provides functionality for renaming environment variable keys
// across one or more .env files in place.
//
// It is useful when refactoring configuration key names to ensure consistency
// across all environment files in a project. Each file is parsed, the target
// key renamed, and the file rewritten preserving all other entries.
//
// Example usage:
//
//	results := renamer.Rename([]string{".env", ".env.staging"}, "DB_HOST", "DATABASE_HOST")
//	for _, r := range results {
//		fmt.Printf("%s: renamed=%v\n", r.File, r.Renamed)
//	}
package renamer
