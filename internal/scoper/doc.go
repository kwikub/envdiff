// Package scoper provides utilities for scoping .env entries to named
// environments such as "production", "staging", or "ci".
//
// Entries are associated with a scope via an inline comment tag:
//
//	DB_URL=postgres://... # @scope:production
//
// The Scope function filters a slice of entries, retaining only those
// that match the requested scope or carry no scope tag at all.
//
// The Tag function annotates specific keys with a scope comment tag,
// replacing any existing scope annotation on those keys.
package scoper
