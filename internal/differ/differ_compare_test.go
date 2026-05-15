package differ

import (
	"testing"

	"github.com/subtlepseudonym/envdiff/internal/parser"
)

func TestCompare_MatchingKeys(t *testing.T) {
	left := []parser.Entry{{Key: "FOO", Value: "bar"}}
	right := []parser.Entry{{Key: "FOO", Value: "bar"}}
	results := Compare(left, right)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "match" {
		t.Errorf("expected match, got %s", results[0].Status)
	}
}

func TestCompare_MismatchedValue(t *testing.T) {
	left := []parser.Entry{{Key: "FOO", Value: "bar"}}
	right := []parser.Entry{{Key: "FOO", Value: "baz"}}
	results := Compare(left, right)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "mismatch" {
		t.Errorf("expected mismatch, got %s", results[0].Status)
	}
	if results[0].Left != "bar" || results[0].Right != "baz" {
		t.Errorf("unexpected values: left=%s right=%s", results[0].Left, results[0].Right)
	}
}

func TestCompare_LeftOnly(t *testing.T) {
	left := []parser.Entry{{Key: "ONLY_LEFT", Value: "x"}}
	right := []parser.Entry{}
	results := Compare(left, right)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "left_only" {
		t.Errorf("expected left_only, got %s", results[0].Status)
	}
}

func TestCompare_RightOnly(t *testing.T) {
	left := []parser.Entry{}
	right := []parser.Entry{{Key: "ONLY_RIGHT", Value: "y"}}
	results := Compare(left, right)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "right_only" {
		t.Errorf("expected right_only, got %s", results[0].Status)
	}
}

func TestCompare_EmptyBothSides(t *testing.T) {
	results := Compare([]parser.Entry{}, []parser.Entry{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
