// Package comparator provides multi-file .env comparison across environments.
package comparator

import (
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// EnvMap maps environment name to its parsed key-value entries.
type EnvMap map[string]map[string]string

// KeyReport summarises how a single key appears across all environments.
type KeyReport struct {
	Key    string
	Values map[string]string // env name -> value ("" means absent)
	Uniform bool             // true if all envs that define the key share the same value
	Missing []string         // env names where the key is absent
}

// Report is the full comparison result across all environments.
type Report struct {
	Environments []string
	Keys         []KeyReport
}

// Compare loads the given files (keyed by environment name) and produces a
// cross-environment comparison report.
func Compare(files map[string]string) (*Report, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("comparator: no files provided")
	}

	envMap := make(EnvMap, len(files))
	for env, path := range files {
		entries, err := parser.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("comparator: parsing %q: %w", path, err)
		}
		m := make(map[string]string, len(entries))
		for _, e := range entries {
			m[e.Key] = e.Value
		}
		envMap[env] = m
	}

	// Collect all unique keys across every environment.
	keySet := make(map[string]struct{})
	for _, m := range envMap {
		for k := range m {
			keySet[k] = struct{}{}
		}
	}

	envNames := make([]string, 0, len(files))
	for env := range files {
		envNames = append(envNames, env)
	}
	sort.Strings(envNames)

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	reports := make([]KeyReport, 0, len(keys))
	for _, key := range keys {
		kr := KeyReport{
			Key:    key,
			Values: make(map[string]string, len(envNames)),
		}
		firstVal := ""
		firstSet := false
		uniform := true
		for _, env := range envNames {
			val, ok := envMap[env][key]
			if !ok {
				kr.Missing = append(kr.Missing, env)
				kr.Values[env] = ""
				uniform = false
				continue
			}
			kr.Values[env] = val
			if !firstSet {
				firstVal = val
				firstSet = true
			} else if val != firstVal {
				uniform = false
			}
		}
		kr.Uniform = uniform
		reports = append(reports, kr)
	}

	return &Report{
		Environments: envNames,
		Keys:         reports,
	}, nil
}
