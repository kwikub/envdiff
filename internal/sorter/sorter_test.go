package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/sorter"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestSort_Alpha_BasicOrder(t *testing.T) {
	input := entries("ZEBRA", "z", "APPLE", "a", "MANGO", "m")
	result := sorter.Sort(input, sorter.SortAlpha)

	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSort_Alpha_DoesNotMutateOriginal(t *testing.T) {
	input := entries("Z", "1", "A", "2")
	originalFirst := input[0].Key

	sorter.Sort(input, sorter.SortAlpha)

	if input[0].Key != originalFirst {
		t.Errorf("original slice was mutated: first key changed to %q", input[0].Key)
	}
}

func TestSort_Group_GroupsKeysByPrefix(t *testing.T) {
	input := entries(
		"APP_NAME", "myapp",
		"DB_HOST", "localhost",
		"APP_PORT", "8080",
		"DB_PASS", "secret",
		"LOG_LEVEL", "info",
	)
	result := sorter.Sort(input, sorter.SortGroup)

	expectedOrder := []string{"APP_NAME", "APP_PORT", "DB_HOST", "DB_PASS", "LOG_LEVEL"}
	for i, e := range result {
		if e.Key != expectedOrder[i] {
			t.Errorf("position %d: got %q, want %q", i, e.Key, expectedOrder[i])
		}
	}
}

func TestSort_Group_NoUnderscore_TreatedAsOwnGroup(t *testing.T) {
	input := entries("ZEBRA", "z", "APP_X", "x", "ALPHA", "a")
	result := sorter.Sort(input, sorter.SortGroup)

	// ALPHA, APP_X, ZEBRA — each is its own group, sorted by group name
	expected := []string{"ALPHA", "APP_X", "ZEBRA"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestSort_EmptyInput(t *testing.T) {
	result := sorter.Sort([]parser.Entry{}, sorter.SortAlpha)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
