package snapshot

import (
	"github.com/user/envdiff/internal/differ"
)

// DriftResult holds the comparison between two snapshots.
type DriftResult struct {
	Baseline *Snapshot
	Current  *Snapshot
	Diffs    []differ.Diff
}

// DetectDrift compares a baseline snapshot against a current snapshot
// and returns any key-level differences (added, removed, modified).
func DetectDrift(baseline, current *Snapshot) (*DriftResult, error) {
	diffs, err := differ.Diff(baseline.Entries, current.Entries)
	if err != nil {
		return nil, err
	}
	return &DriftResult{
		Baseline: baseline,
		Current:  current,
		Diffs:    diffs,
	}, nil
}

// HasDrift returns true when there is at least one non-unchanged diff.
func (d *DriftResult) HasDrift() bool {
	for _, diff := range d.Diffs {
		if diff.Status != differ.StatusUnchanged {
			return true
		}
	}
	return false
}

// Summary returns counts of added, removed, and modified keys.
func (d *DriftResult) Summary() (added, removed, modified int) {
	for _, diff := range d.Diffs {
		switch diff.Status {
		case differ.StatusAdded:
			added++
		case differ.StatusRemoved:
			removed++
		case differ.StatusModified:
			modified++
		}
	}
	return
}
