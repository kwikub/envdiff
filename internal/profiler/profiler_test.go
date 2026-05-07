package profiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/profiler"
)

func makeEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_ENV", Value: "staging"},
		{Key: "DB_HOST", Value: "db.staging.internal"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := profiler.Profile{Name: "staging", Entries: makeEntries()}

	if err := profiler.Save(dir, p); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := profiler.Load(dir, "staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Name != p.Name {
		t.Errorf("name: got %q, want %q", loaded.Name, p.Name)
	}
	if len(loaded.Entries) != len(p.Entries) {
		t.Fatalf("entries len: got %d, want %d", len(loaded.Entries), len(p.Entries))
	}
	for i, e := range loaded.Entries {
		if e.Key != p.Entries[i].Key || e.Value != p.Entries[i].Value {
			t.Errorf("entry %d: got %+v, want %+v", i, e, p.Entries[i])
		}
	}
}

func TestSave_EmptyName_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	err := profiler.Save(dir, profiler.Profile{Name: "", Entries: makeEntries()})
	if err == nil {
		t.Fatal("expected error for empty profile name, got nil")
	}
}

func TestLoad_NotFound_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	_, err := profiler.Load(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error loading nonexistent profile, got nil")
	}
}

func TestList_ReturnsProfileNames(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"dev", "staging", "prod"} {
		p := profiler.Profile{Name: name, Entries: makeEntries()}
		if err := profiler.Save(dir, p); err != nil {
			t.Fatalf("Save %q: %v", name, err)
		}
	}
	// add a non-json file to ensure it is ignored
	_ = os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("ignore me"), 0o644)

	names, err := profiler.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("got %d profiles, want 3: %v", len(names), names)
	}
}

func TestList_EmptyDir_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	names, err := profiler.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}
