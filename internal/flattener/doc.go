// Package flattener provides Flatten, which merges multiple slices of
// parser.Entry values into a single deduplicated slice.
//
// Conflict resolution is governed by a Strategy:
//
//   - StrategyFirst  – the first value seen for a key wins.
//   - StrategyLast   – the last value seen for a key wins.
//   - StrategyError  – any duplicate key causes an immediate error.
//
// Flatten preserves the insertion order of keys as they are first encountered
// across the provided groups, making the output deterministic.
package flattener
