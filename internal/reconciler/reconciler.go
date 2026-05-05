package reconciler

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

// Action represents what should be done with a key during reconciliation.
type Action string

const (
	ActionAdd    Action = "add"
	ActionRemove Action = "remove"
	ActionUpdate Action = "update"
	ActionKeep   Action = "keep"
)

// Suggestion represents a reconciliation suggestion for a single key.
type Suggestion struct {
	Key    string
	Action Action
	Value  string
	Reason string
}

// Reconcile compares a source env file against a target env file and returns
// a list of suggestions to bring the target in sync with the source.
func Reconcile(sourcePath, targetPath string) ([]Suggestion, error) {
	sourceEntries, err := parser.Parse(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("parsing source: %w", err)
	}

	targetEntries, err := parser.Parse(targetPath)
	if err != nil {
		return nil, fmt.Errorf("parsing target: %w", err)
	}

	diffs := differ.Diff(sourceEntries, targetEntries)

	var suggestions []Suggestion
	for _, d := range diffs {
		switch d.Status {
		case differ.Added:
			suggestions = append(suggestions, Suggestion{
				Key:    d.Key,
				Action: ActionRemove,
				Value:  d.TargetValue,
				Reason: "key exists in target but not in source",
			})
		case differ.Removed:
			suggestions = append(suggestions, Suggestion{
				Key:    d.Key,
				Action: ActionAdd,
				Value:  d.SourceValue,
				Reason: "key missing from target",
			})
		case differ.Modified:
			suggestions = append(suggestions, Suggestion{
				Key:    d.Key,
				Action: ActionUpdate,
				Value:  d.SourceValue,
				Reason: "value differs from source",
			})
		case differ.Unchanged:
			suggestions = append(suggestions, Suggestion{
				Key:    d.Key,
				Action: ActionKeep,
				Value:  d.SourceValue,
				Reason: "in sync",
			})
		}
	}
	return suggestions, nil
}

// Format renders suggestions as human-readable lines.
func Format(suggestions []Suggestion, maskSecrets bool) string {
	var sb strings.Builder
	for _, s := range suggestions {
		val := s.Value
		if maskSecrets && differ.IsSecret(s.Key) {
			val = "***"
		}
		sb.WriteString(fmt.Sprintf("[%s] %s=%s  # %s\n", s.Action, s.Key, val, s.Reason))
	}
	return sb.String()
}
