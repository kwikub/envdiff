// Package flattener collapses multiple env entry slices into a single
// deduplicated slice, resolving conflicts according to a chosen strategy.
package flattener

import (
	"errors"
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Strategy controls how duplicate keys are resolved during flattening.
type Strategy string

const (
	// StrategyFirst keeps the first occurrence of a key.
	StrategyFirst Strategy = "first"
	// StrategyLast keeps the last occurrence of a key.
	StrategyLast Strategy = "last"
	// StrategyError returns an error if any duplicate key is encountered.
	StrategyError Strategy = "error"
)

// ErrDuplicateKey is returned by Flatten when strategy is StrategyError and
// a duplicate key is detected.
var ErrDuplicateKey = errors.New("duplicate key")

// Flatten merges the provided entry slices into one, applying strategy to
// resolve conflicts. The relative order of first-seen keys is preserved.
func Flatten(strategy Strategy, groups ...[]parser.Entry) ([]parser.Entry, error) {
	seen := make(map[string]int) // key -> index in result
	var result []parser.Entry

	for _, entries := range groups {
		for _, e := range entries {
			idx, exists := seen[e.Key]
			switch {
			case !exists:
				seen[e.Key] = len(result)
				result = append(result, e)
			case strategy == StrategyFirst:
				// keep existing — no-op
			case strategy == StrategyLast:
				result[idx] = e
			case strategy == StrategyError:
				return nil, fmt.Errorf("%w: %q", ErrDuplicateKey, e.Key)
			default:
				return nil, fmt.Errorf("unknown strategy %q", strategy)
			}
		}
	}

	return result, nil
}
