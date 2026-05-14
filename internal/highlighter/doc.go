// Package highlighter renders colourised, human-readable output for .env
// file diffs produced by the differ package.
//
// It accepts a slice of differ.Result values and formats each entry with
// ANSI colour codes:
//
//   - Green  (+) for added keys
//   - Red    (-) for removed keys
//   - Yellow (~) for modified keys (shown as a red/green pair)
//   - Gray       for unchanged keys (hidden by default)
//
// Colour output can be suppressed via Options.NoColor for use in CI
// pipelines or when redirecting to a file.
package highlighter
