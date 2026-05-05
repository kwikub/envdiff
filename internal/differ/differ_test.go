package differ_test

import (
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestDiff_AddedKey(t *testing.T) {
	base := entries("HOST", "localhost")
	other := entries("HOST", "localhost", "PORT", "8080")

	result := differ.Diff(base, other)
	found := findByKey(result, "PORT")
	if found == nil || found.Type != differ.Added {
		t.Errorf("expected PORT to be Added, got %+v", found)
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := entries("HOST", "localhost", "PORT", "8080")
	other := entries("HOST", "localhost")

	result := differ.Diff(base, other)
	found := findByKey(result, "PORT")
	if found == nil || found.Type != differ.Removed {
		t.Errorf("expected PORT to be Removed, got %+v", found)
	}
}

func TestDiff_ModifiedKey(t *testing.T) {
	base := entries("HOST", "localhost")
	other := entries("HOST", "prod.example.com")

	result := differ.Diff(base, other)
	found := findByKey(result, "HOST")
	if found == nil || found.Type != differ.Modified {
		t.Errorf("expected HOST to be Modified, got %+v", found)
	}
	if found.BaseVal != "localhost" || found.OtherVal != "prod.example.com" {
		t.Errorf("unexpected values: %+v", found)
	}
}

func TestDiff_UnchangedKey(t *testing.T) {
	base := entries("HOST", "localhost")
	other := entries("HOST", "localhost")

	result := differ.Diff(base, other)
	found := findByKey(result, "HOST")
	if found == nil || found.Type != differ.Unchanged {
		t.Errorf("expected HOST to be Unchanged, got %+v", found)
	}
}

func TestDiff_EmptyFiles(t *testing.T) {
	result := differ.Diff(nil, nil)
	if len(result) != 0 {
		t.Errorf("expected empty diff, got %d entries", len(result))
	}
}

func findByKey(entries []differ.DiffEntry, key string) *differ.DiffEntry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}
