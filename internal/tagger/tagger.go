// Package tagger provides functionality for tagging .env entries with
// arbitrary labels, enabling grouping, filtering, and annotation workflows.
package tagger

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Tag represents a label attached to one or more keys.
type Tag struct {
	Name string
	Keys []string
}

// Result holds the tagged output for a single entry.
type Result struct {
	Key   string
	Value string
	Tags  []string
}

// Options controls tagging behaviour.
type Options struct {
	// TagSecrets automatically applies the "secret" tag to keys that look like secrets.
	TagSecrets bool
	// TagEmpty automatically applies the "empty" tag to keys with blank values.
	TagEmpty bool
	// Extra maps tag names to explicit key lists.
	Extra map[string][]string
}

// Apply tags entries according to the supplied options and returns a slice of
// Result values preserving the original order.
func Apply(entries []parser.Entry, opts Options) ([]Result, error) {
	if entries == nil {
		return nil, fmt.Errorf("tagger: entries must not be nil")
	}

	// Build a reverse lookup: key → set of tags.
	keyTags := make(map[string]map[string]struct{})
	for _, e := range entries {
		keyTags[e.Key] = make(map[string]struct{})
	}

	if opts.TagSecrets {
		for _, e := range entries {
			if isSecret(e.Key) {
				keyTags[e.Key]["secret"] = struct{}{}
			}
		}
	}

	if opts.TagEmpty {
		for _, e := range entries {
			if strings.TrimSpace(e.Value) == "" {
				keyTags[e.Key]["empty"] = struct{}{}
			}
		}
	}

	for tag, keys := range opts.Extra {
		for _, k := range keys {
			if _, ok := keyTags[k]; ok {
				keyTags[k][tag] = struct{}{}
			}
		}
	}

	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		tags := make([]string, 0, len(keyTags[e.Key]))
		for t := range keyTags[e.Key] {
			tags = append(tags, t)
		}
		sort.Strings(tags)
		results = append(results, Result{Key: e.Key, Value: e.Value, Tags: tags})
	}
	return results, nil
}

// FilterByTag returns only those results that carry the given tag.
func FilterByTag(results []Result, tag string) []Result {
	var out []Result
	for _, r := range results {
		for _, t := range r.Tags {
			if t == tag {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

var secretSuffixes = []string{"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CERT", "CREDENTIAL"}

func isSecret(key string) bool {
	upper := strings.ToUpper(key)
	for _, s := range secretSuffixes {
		if strings.Contains(upper, s) {
			return true
		}
	}
	return false
}
