// Package scanner detects potentially sensitive or misconfigured values
// across .env files by applying heuristic rules.
package scanner

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Severity represents the severity level of a scan finding.
type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
)

// Finding describes a single issue detected during scanning.
type Finding struct {
	Key      string
	Value    string
	Rule     string
	Message  string
	Severity Severity
}

var (
	highEntropyPattern = regexp.MustCompile(`[A-Za-z0-9+/]{32,}`)
	urlWithCredentials = regexp.MustCompile(`[a-z]+://[^:]+:[^@]+@`)
	privateKeyHeader   = regexp.MustCompile(`-----BEGIN .* PRIVATE KEY-----`)
)

// Scan parses the given file and applies heuristic rules to each entry,
// returning a list of findings.
func Scan(path string) ([]Finding, error) {
	entries, err := parser.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("scanner: parse %q: %w", path, err)
	}

	var findings []Finding
	for _, e := range entries {
		findings = append(findings, applyRules(e.Key, e.Value)...)
	}
	return findings, nil
}

func applyRules(key, value string) []Finding {
	var out []Finding

	if privateKeyHeader.MatchString(value) {
		out = append(out, Finding{
			Key:      key,
			Value:    "***",
			Rule:     "private-key-header",
			Message:  "value appears to contain a PEM private key",
			Severity: SeverityError,
		})
	}

	if urlWithCredentials.MatchString(value) {
		out = append(out, Finding{
			Key:      key,
			Value:    "***",
			Rule:     "url-credentials",
			Message:  "value contains a URL with embedded credentials",
			Severity: SeverityError,
		})
	}

	if len(value) >= 32 && highEntropyPattern.MatchString(value) && !strings.Contains(value, " ") {
		out = append(out, Finding{
			Key:      key,
			Value:    "***",
			Rule:     "high-entropy",
			Message:  "value has high entropy and may be a secret token",
			Severity: SeverityWarning,
		})
	}

	if value == "changeme" || value == "password" || value == "secret" {
		out = append(out, Finding{
			Key:      key,
			Value:    value,
			Rule:     "placeholder-value",
			Message:  "value looks like an unset placeholder",
			Severity: SeverityInfo,
		})
	}

	return out
}
