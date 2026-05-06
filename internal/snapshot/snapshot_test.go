package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/snapshot"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestTake_ReturnsSnapshot(t *testing.T) {
	path := writeTempEnv(t, "APP=hello\nDB=postgres\n")
	snap, err := snapshot.Take(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Source != path {
		t.Errorf("expected source %q, got %q", path, snap.Source)
	}
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	snap, err := snapshot.Take(path)
	if err != nil {
		t.Fatalf("take: %v", err)
	}

	dest := filepath.Join(t.TempDir(), "snap.json")
	if err := snapshot.Save(snap, dest); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Source != snap.Source {
		t.Errorf("source mismatch: got %q", loaded.Source)
	}
	if len(loaded.Entries) != len(snap.Entries) {
		t.Errorf("entries mismatch: got %d", len(loaded.Entries))
	}
}

func TestToMap_ConvertsEntries(t *testing.T) {
	path := writeTempEnv(t, "A=1\nB=2\n")
	snap, _ := snapshot.Take(path)
	m := snap.ToMap()
	if m["A"] != "1" || m["B"] != "2" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestLoad_InvalidFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
