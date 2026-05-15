package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/subtlepseudonym/envdiff/internal/differ"
	"github.com/subtlepseudonym/envdiff/internal/parser"
)

// runCompare implements the `envdiff compare` sub-command.
// Usage: envdiff compare [--json] [--mask] <file-a> <file-b>
func runCompare(args []string, jsonOut bool, maskSecrets bool) error {
	if len(args) < 2 {
		return fmt.Errorf("compare requires two file paths")
	}

	left, err := parser.Parse(args[0])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[0], err)
	}
	right, err := parser.Parse(args[1])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", args[1], err)
	}

	results := differ.Compare(left, right)
	if maskSecrets {
		results = differ.MaskSecrets(results)
	}

	if jsonOut {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	for _, r := range results {
		switch r.Status {
		case "match":
			fmt.Printf("  %s=%s\n", r.Key, r.Left)
		case "mismatch":
			fmt.Printf("~ %s: %s → %s\n", r.Key, r.Left, r.Right)
		case "left_only":
			fmt.Printf("- %s=%s\n", r.Key, r.Left)
		case "right_only":
			fmt.Printf("+ %s=%s\n", r.Key, r.Right)
		}
	}
	return nil
}
