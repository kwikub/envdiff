// Package auditor provides functionality to audit .env files for
// sensitive key exposure, unused variables, and policy violations.
package auditor

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

// Severity represents the level of an audit finding.
type Severity string

const (
	SeverityInfo    Severity = "INFO"
	SeverityWarning Severity = "WARNING"
	SeverityError   Severity = "ERROR"
)

// Finding represents a single audit result.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Result holds the complete output of an audit run.
type Result struct {
	Findings []Finding
	HasError bool
}

// Audit parses the given .env file and checks for policy issues.
func Audit(path string) (*Result, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("audit: failed to parse %q: %w", path, err)
	}

	result := &Result{}

	for _, e := range entries {
		// Check for secret keys with empty values
		if differ.IsSecret(e.Key) && strings.TrimSpace(e.Value) == "" {
			result.add(Finding{
				Key:      e.Key,
				Message:  "secret key has an empty value",
				Severity: SeverityError,
			})
		}

		// Warn about keys with whitespace in values
		if strings.Contains(e.Value, " ") && !strings.HasPrefix(e.Value, "\"") {
			result.add(Finding{
				Key:      e.Key,
				Message:  "value contains unquoted whitespace",
				Severity: SeverityWarning,
			})
		}

		// Info: plaintext secret detected (non-empty secret key)
		if differ.IsSecret(e.Key) && e.Value != "" {
			result.add(Finding{
				Key:      e.Key,
				Message:  "plaintext secret value present in file",
				Severity: SeverityInfo,
			})
		}
	}

	return result, nil
}

func (r *Result) add(f Finding) {
	if f.Severity == SeverityError {
		r.HasError = true
	}
	r.Findings = append(r.Findings, f)
}
