// Package validator provides functionality to validate .env files
// against a reference template, checking for required keys and value formats.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// ValidationError represents a single validation issue found in an env file.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// Result holds the outcome of a validation run.
type Result struct {
	Valid  bool
	Errors []ValidationError
}

// Rule defines a validation rule applied to a specific key.
type Rule struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp // optional regex the value must match
}

// Validate checks the entries parsed from envFile against the provided rules.
// It also accepts a template file path; any key present in the template but
// missing from envFile is reported as an error when no explicit rule overrides.
func Validate(envFile string, templateFile string, rules []Rule) (Result, error) {
	envEntries, err := parser.Parse(envFile)
	if err != nil {
		return Result{}, fmt.Errorf("parsing env file: %w", err)
	}

	tmplEntries, err := parser.Parse(templateFile)
	if err != nil {
		return Result{}, fmt.Errorf("parsing template file: %w", err)
	}

	envMap := toMap(envEntries)
	tmplMap := toMap(tmplEntries)
	ruleMap := make(map[string]Rule, len(rules))
	for _, r := range rules {
		ruleMap[strings.ToUpper(r.Key)] = r
	}

	var errs []ValidationError

	// Check all template keys exist in env file.
	for key := range tmplMap {
		if rule, ok := ruleMap[key]; ok && !rule.Required {
			continue
		}
		if _, present := envMap[key]; !present {
			errs = append(errs, ValidationError{Key: key, Message: "required key is missing"})
		}
	}

	// Apply pattern rules.
	for _, rule := range rules {
		key := strings.ToUpper(rule.Key)
		val, exists := envMap[key]
		if rule.Required && !exists {
			// Already captured above if it was in template; add if not.
			if _, inTmpl := tmplMap[key]; !inTmpl {
				errs = append(errs, ValidationError{Key: key, Message: "required key is missing"})
			}
			continue
		}
		if rule.Pattern != nil && exists && !rule.Pattern.MatchString(val) {
			errs = append(errs, ValidationError{
				Key:     key,
				Message: fmt.Sprintf("value does not match required pattern %s", rule.Pattern.String()),
			})
		}
	}

	return Result{Valid: len(errs) == 0, Errors: errs}, nil
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
