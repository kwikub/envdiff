package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/interpolator"
	"github.com/user/envdiff/internal/parser"
)

// runInterpolate is the handler for the `interpolate` sub-command.
// Usage: envdiff interpolate <file> [--strict] [--os] [--override KEY=VALUE ...]
func runInterpolate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff interpolate <file> [--strict] [--os] [--override KEY=VALUE]")
	}

	filePath := args[0]
	var strict, fallbackOS bool
	overrides := map[string]string{}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--strict":
			strict = true
		case "--os":
			fallbackOS = true
		case "--override":
			i++
			if i >= len(args) {
				return fmt.Errorf("--override requires KEY=VALUE argument")
			}
			parts := strings.SplitN(args[i], "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid override %q: expected KEY=VALUE", args[i])
			}
			overrides[parts[0]] = parts[1]
		}
	}

	entries, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("parse %q: %w", filePath, err)
	}

	opts := interpolator.Options{Strict: strict, FallbackToOS: fallbackOS}
	resolved, err := interpolator.Interpolate(entries, overrides, opts)
	if err != nil {
		return fmt.Errorf("interpolation failed: %w", err)
	}

	out := make(map[string]string, len(resolved))
	for _, e := range resolved {
		out[e.Key] = e.Value
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encode output: %w", err)
	}
	return nil
}
