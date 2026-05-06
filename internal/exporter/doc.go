// Package exporter converts parsed .env entries into various output formats
// suitable for different deployment targets.
//
// Supported formats:
//
//   - env    — standard KEY=VALUE format (values with spaces are quoted)
//   - shell  — KEY=VALUE prefixed with 'export' for sourcing in bash/zsh
//   - docker — Docker-compatible env-file format (no quoting)
//   - json   — JSON object mapping keys to string values
//
// Example usage:
//
//	exporter.Export(entries, exporter.FormatShell, os.Stdout)
package exporter
