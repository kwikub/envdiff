package differ

import "github.com/subtlepseudonym/envdiff/internal/parser"

// MaskSecrets returns a copy of results with secret values replaced by "***".
// It operates on CompareResult slices produced by Compare so that callers can
// safely render output without leaking credentials.
func MaskSecrets(results []CompareResult) []CompareResult {
	out := make([]CompareResult, len(results))
	for i, r := range results {
		if IsSecret(parser.Entry{Key: r.Key}) {
			r.Left = "***"
			r.Right = "***"
		}
		out[i] = r
	}
	return out
}
