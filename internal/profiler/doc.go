// Package profiler manages named environment profiles for envdiff.
//
// A profile is a named snapshot of environment entries that can be saved to
// disk and loaded later. Profiles are stored as JSON files under a configurable
// directory (e.g. ~/.config/envdiff/profiles/).
//
// Typical usage:
//
//	// Save current env as a profile
//	p := profiler.Profile{Name: "staging", Entries: entries}
//	profiler.Save("/path/to/profiles", p)
//
//	// Load a profile by name
//	p, err := profiler.Load("/path/to/profiles", "staging")
//
//	// List available profiles
//	names, err := profiler.List("/path/to/profiles")
package profiler
