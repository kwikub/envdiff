package flattener_test

import (
	"errors"
	"testing"

	"github.com/user/envdiff/internal/flattener"
	"github.com/user/envdiff/internal/parser"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestFlatten_NoGroups_ReturnsEmpty(t *testing.T) {
	result, err := flattener.Flatten(flattener.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(result))
	}
}

func TestFlatten_SingleGroup_NoDuplicates(t *testing.T) {
	result, err := flattener.Flatten(flattener.StrategyFirst,
		entries("A", "1", "B", "2"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFlatten_StrategyFirst_KeepsFirstValue(t *testing.T) {
	result, err := flattener.Flatten(flattener.StrategyFirst,
		entries("KEY", "original"),
		entries("KEY", "override"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Value != "original" {
		t.Errorf("expected %q, got %q", "original", result[0].Value)
	}
}

func TestFlatten_StrategyLast_KeepsLastValue(t *testing.T) {
	result, err := flattener.Flatten(flattener.StrategyLast,
		entries("KEY", "original"),
		entries("KEY", "override"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Value != "override" {
		t.Errorf("expected %q, got %q", "override", result[0].Value)
	}
}

func TestFlatten_StrategyError_ReturnsDuplicateKeyError(t *testing.T) {
	_, err := flattener.Flatten(flattener.StrategyError,
		entries("KEY", "a"),
		entries("KEY", "b"),
	)
	if !errors.Is(err, flattener.ErrDuplicateKey) {
		t.Errorf("expected ErrDuplicateKey, got %v", err)
	}
}

func TestFlatten_StrategyError_NoDuplicates_ReturnsNoError(t *testing.T) {
	_, err := flattener.Flatten(flattener.StrategyError,
		entries("A", "1"),
		entries("B", "2"),
	)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestFlatten_PreservesOrderOfFirstSeen(t *testing.T) {
	result, err := flattener.Flatten(flattener.StrategyFirst,
		entries("C", "3", "A", "1"),
		entries("B", "2", "A", "99"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"C", "A", "B"}
	for i, e := range result {
		if e.Key != want[i] {
			t.Errorf("position %d: expected key %q, got %q", i, want[i], e.Key)
		}
	}
}
