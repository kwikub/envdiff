package renamer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/renamer"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestRename_KeyExists_RenamedSuccessfully(t *testing.T) {
	p := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")

	results := renamer.Rename([]string{p}, "DB_HOST", "DATABASE_HOST")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Renamed {
		t.Errorf("expected Renamed=true")
	}
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}

	entries, err := parser.Parse(p)
	if err != nil {
		t.Fatalf("parse after rename: %v", err)
	}
	for _, e := range entries {
		if e.Key == "DB_HOST" {
			t.Errorf("old key DB_HOST still present")
		}
	}
	found := false
	for _, e := range entries {
		if e.Key == "DATABASE_HOST" {
			found = true
			if e.Value != "localhost" {
				t.Errorf("expected value localhost, got %q", e.Value)
			}
		}
	}
	if !found {
		t.Errorf("new key DATABASE_HOST not found")
	}
}

func TestRename_KeyMissing_NotRenamed(t *testing.T) {
	p := writeTempEnv(t, "APP_ENV=production\n")

	results := renamer.Rename([]string{p}, "MISSING_KEY", "NEW_KEY")
	if results[0].Renamed {
		t.Errorf("expected Renamed=false for missing key")
	}
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}
}

func TestRename_MultipleFiles(t *testing.T) {
	p1 := writeTempEnv(t, "SECRET_KEY=abc\n")
	p2 := writeTempEnv(t, "OTHER=xyz\n")

	results := renamer.Rename([]string{p1, p2}, "SECRET_KEY", "APP_SECRET")
	if len(results) != 2 {
		t.Fatalf("expected 2 results")
	}
	if !results[0].Renamed {
		t.Errorf("expected file1 to be renamed")
	}
	if results[1].Renamed {
		t.Errorf("expected file2 not to be renamed")
	}
}

func TestRename_InvalidFile_ReturnsError(t *testing.T) {
	results := renamer.Rename([]string{"/nonexistent/.env"}, "KEY", "NEW_KEY")
	if results[0].Err == nil {
		t.Errorf("expected error for nonexistent file")
	}
}
