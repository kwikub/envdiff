// Package linter provides static analysis checks for .env files,
// warning about common issues such as duplicate keys, empty values,
// and keys that do not follow naming conventions.
package linter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Severity represents the level of a lint finding.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Finding describes a single lint issue found in a .env file.
type Finding struct {
	Line     int
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] line %d: %s — %s", f.Severity, f.Line, f.Key, f.Message)
}

// validKeyRe matches conventional env var names: uppercase letters, digits, underscores.
var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Lint parses the file at path and returns a slice of Findings.
func Lint(path string) ([]Finding, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("linter: parse %q: %w", path, err)
	}

	var findings []Finding
	seen := make(map[string]int) // key -> first line number (1-based index)

	for i, e := range entries {
		lineNum := i + 1

		// Duplicate key check
		if prev, ok := seen[e.Key]; ok {
			findings = append(findings, Finding{
				Line:     lineNum,
				Key:      e.Key,
				Message:  fmt.Sprintf("duplicate key (first seen at line %d)", prev),
				Severity: SeverityError,
			})
		} else {
			seen[e.Key] = lineNum
		}

		// Empty value check
		if strings.TrimSpace(e.Value) == "" {
			findings = append(findings, Finding{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "empty value",
				Severity: SeverityWarning,
			})
		}

		// Naming convention check
		if !validKeyRe.MatchString(e.Key) {
			findings = append(findings, Finding{
				Line:     lineNum,
				Key:      e.Key,
				Message:  "key does not follow UPPER_SNAKE_CASE convention",
				Severity: SeverityWarning,
			})
		}
	}

	return findings, nil
}

// HasErrors reports whether any of the provided findings have SeverityError.
// This is useful for callers that need to decide whether to treat lint results
// as a hard failure (e.g. in CI pipelines).
func HasErrors(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == SeverityError {
			return true
		}
	}
	return false
}
