// Package pinner implements env lockfile functionality for envdiff.
//
// It allows users to "pin" the current values of a .env file into a JSON
// lockfile. The lockfile can later be used to verify that the .env file
// has not drifted from the pinned state — useful for CI checks and
// auditing environment consistency across deployments.
//
// Usage:
//
//	// Create a lockfile from the current .env
//	pf, err := pinner.Pin(".env", ".env.lock")
//
//	// Verify the .env has not changed since it was pinned
//	drifted, err := pinner.Verify(".env", ".env.lock")
package pinner
