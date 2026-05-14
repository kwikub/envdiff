package differ

import (
	"testing"
)

func TestCompare_MatchingKeys(t *testing.T) {
	left := entries("A=1", "B=2")
	right := entries("A=1", "B=2")

	results := Compare(left, right)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != "match" {
			t.Errorf("key %s: expected match, got %s", r.Key, r.Status)
		}
	}
}

func TestCompare_MismatchedValue(t *testing.T) {
	left := entries("A=1")
	right := entries("A=2")

	results := Compare(left, right)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "mismatch" {
		t.Errorf("expected mismatch, got %s", results[0].Status)
	}
	if results[0].Left != "1" || results[0].Right != "2" {
		t.Errorf("unexpected values: left=%s right=%s", results[0].Left, results[0].Right)
	}
}

func TestCompare_LeftOnly(t *testing.T) {
	left := entries("A=1", "B=2")
	right := entries("A=1")

	results := Compare(left, right)
	var found *CompareResult
	for i := range results {
		if results[i].Key == "B" {
			found = &results[i]
		}
	}
	if found == nil {
		t.Fatal("expected result for key B")
	}
	if found.Status != "left_only" {
		t.Errorf("expected left_only, got %s", found.Status)
	}
	if found.Right != "" {
		t.Errorf("expected empty right value")
	}
}

func TestCompare_RightOnly(t *testing.T) {
	left := entries("A=1")
	right := entries("A=1", "C=3")

	results := Compare(left, right)
	var found *CompareResult
	for i := range results {
		if results[i].Key == "C" {
			found = &results[i]
		}
	}
	if found == nil {
		t.Fatal("expected result for key C")
	}
	if found.Status != "right_only" {
		t.Errorf("expected right_only, got %s", found.Status)
	}
	if found.Left != "" {
		t.Errorf("expected empty left value")
	}
}

func TestCompare_EmptyBothSides(t *testing.T) {
	results := Compare([]Entry{}, []Entry{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
