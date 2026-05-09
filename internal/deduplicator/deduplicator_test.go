package deduplicator_test

import (
	"testing"

	"github.com/user/envdiff/internal/deduplicator"
	"github.com/user/envdiff/internal/parser"
)

func entries(kvs ...string) []parser.Entry {
	out := make([]parser.Entry, 0, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1], Line: i/2 + 1})
	}
	return out
}

func TestDeduplicate_NoDuplicates(t *testing.T) {
	in := entries("A", "1", "B", "2", "C", "3")
	res, err := deduplicator.Deduplicate(in, deduplicator.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.Entries))
	}
	if len(res.Removed) != 0 {
		t.Errorf("expected no removed entries, got %d", len(res.Removed))
	}
}

func TestDeduplicate_StrategyFirst_KeepsFirst(t *testing.T) {
	in := entries("A", "first", "B", "only", "A", "second")
	res, err := deduplicator.Deduplicate(in, deduplicator.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "first" {
		t.Errorf("expected value %q, got %q", "first", res.Entries[0].Value)
	}
	if len(res.Removed) != 1 || res.Removed[0].Value != "second" {
		t.Errorf("expected removed entry with value %q", "second")
	}
}

func TestDeduplicate_StrategyLast_KeepsLast(t *testing.T) {
	in := entries("A", "first", "B", "only", "A", "second")
	res, err := deduplicator.Deduplicate(in, deduplicator.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "second" {
		t.Errorf("expected value %q, got %q", "second", res.Entries[0].Value)
	}
	if len(res.Removed) != 1 || res.Removed[0].Value != "first" {
		t.Errorf("expected removed entry with value %q", "first")
	}
}

func TestDeduplicate_UnknownStrategy_ReturnsError(t *testing.T) {
	in := entries("A", "1")
	_, err := deduplicator.Deduplicate(in, deduplicator.Strategy("unknown"))
	if err == nil {
		t.Fatal("expected error for unknown strategy, got nil")
	}
}

func TestDeduplicate_MultipleOccurrences_StrategyLast(t *testing.T) {
	in := entries("X", "v1", "X", "v2", "X", "v3")
	res, err := deduplicator.Deduplicate(in, deduplicator.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "v3" {
		t.Errorf("expected value %q, got %q", "v3", res.Entries[0].Value)
	}
	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed entries, got %d", len(res.Removed))
	}
}
