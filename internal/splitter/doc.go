// Package splitter partitions a flat .env file into multiple smaller files
// based on a configurable strategy.
//
// Three strategies are supported:
//
//   - StrategyPrefix (default): groups keys by their underscore-delimited
//     prefix, writing each group to <prefix>.env.
//
//   - StrategyAlpha: distributes keys into N evenly-sized buckets named
//     part00.env, part01.env, … using the first character of each key.
//
//   - Custom mapping: an explicit map[outputFile][]keyPrefix lets callers
//     route specific prefixes to named output files; unmatched keys fall
//     through to other.env.
//
// Example:
//
//	counts, err := splitter.Split(".env", "./split", splitter.Options{
//		Strategy: splitter.StrategyPrefix,
//	})
package splitter
