package snapshot_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/snapshot"
)

func makeSnap(entries []parser.Entry) *snapshot.Snapshot {
	return &snapshot.Snapshot{Source: "test", Entries: entries}
}

func TestDetectDrift_NoDrift(t *testing.T) {
	entries := []parser.Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	result, err := snapshot.DetectDrift(makeSnap(entries), makeSnap(entries))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestDetectDrift_AddedKey(t *testing.T) {
	baseline := []parser.Entry{{Key: "A", Value: "1"}}
	current := []parser.Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	result, err := snapshot.DetectDrift(makeSnap(baseline), makeSnap(current))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDrift() {
		t.Error("expected drift")
	}
	added, removed, modified := result.Summary()
	if added != 1 || removed != 0 || modified != 0 {
		t.Errorf("unexpected summary: added=%d removed=%d modified=%d", added, removed, modified)
	}
}

func TestDetectDrift_ModifiedKey(t *testing.T) {
	baseline := []parser.Entry{{Key: "A", Value: "old"}}
	current := []parser.Entry{{Key: "A", Value: "new"}}
	result, err := snapshot.DetectDrift(makeSnap(baseline), makeSnap(current))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, _, modified := result.Summary()
	if modified != 1 {
		t.Errorf("expected 1 modified, got %d", modified)
	}
}

func TestDetectDrift_RemovedKey(t *testing.T) {
	baseline := []parser.Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	current := []parser.Entry{{Key: "A", Value: "1"}}
	result, err := snapshot.DetectDrift(makeSnap(baseline), makeSnap(current))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, removed, _ := result.Summary()
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}
}

func TestDriftResult_DiffsAccessible(t *testing.T) {
	baseline := []parser.Entry{{Key: "X", Value: "a"}}
	current := []parser.Entry{{Key: "X", Value: "b"}}
	result, _ := snapshot.DetectDrift(makeSnap(baseline), makeSnap(current))
	for _, d := range result.Diffs {
		if d.Key == "X" && d.Status != differ.StatusModified {
			t.Errorf("expected modified status for X")
		}
	}
}

func TestDetectDrift_NilSnapshot(t *testing.T) {
	entries := []parser.Entry{{Key: "A", Value: "1"}}
	_, err := snapshot.DetectDrift(nil, makeSnap(entries))
	if err == nil {
		t.Error("expected error when baseline snapshot is nil")
	}
	_, err = snapshot.DetectDrift(makeSnap(entries), nil)
	if err == nil {
		t.Error("expected error when current snapshot is nil")
	}
}
