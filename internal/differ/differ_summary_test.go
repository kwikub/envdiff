package differ

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func TestSummarise_AllStatuses(t *testing.T) {
	diffs := []DiffEntry{
		{Key: "A", Status: StatusAdded},
		{Key: "B", Status: StatusAdded},
		{Key: "C", Status: StatusRemoved},
		{Key: "D", Status: StatusModified},
		{Key: "E", Status: StatusUnchanged},
		{Key: "F", Status: StatusUnchanged},
	}
	s := Summarise(diffs)
	if s.Added != 2 {
		t.Errorf("Added: want 2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("Modified: want 1, got %d", s.Modified)
	}
	if s.Unchanged != 2 {
		t.Errorf("Unchanged: want 2, got %d", s.Unchanged)
	}
	if s.Total != 6 {
		t.Errorf("Total: want 6, got %d", s.Total)
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := Summarise(nil)
	if s.Total != 0 {
		t.Errorf("Total: want 0, got %d", s.Total)
	}
}

func TestSummarise_OnlyUnchanged(t *testing.T) {
	diffs := []DiffEntry{
		{Key: "X", Status: StatusUnchanged},
		{Key: "Y", Status: StatusUnchanged},
	}
	s := Summarise(diffs)
	if s.Added != 0 || s.Removed != 0 || s.Modified != 0 {
		t.Errorf("expected no changes, got added=%d removed=%d modified=%d", s.Added, s.Removed, s.Modified)
	}
	if s.Unchanged != 2 {
		t.Errorf("Unchanged: want 2, got %d", s.Unchanged)
	}
}

func TestSummarise_MatchesDiffOutput(t *testing.T) {
	base := []parser.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	target := []parser.Entry{
		{Key: "FOO", Value: "changed"},
		{Key: "NEW", Value: "value"},
	}
	diffs := Diff(base, target)
	s := Summarise(diffs)
	if s.Added != 1 {
		t.Errorf("Added: want 1, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("Modified: want 1, got %d", s.Modified)
	}
}
