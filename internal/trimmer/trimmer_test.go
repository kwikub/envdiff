package trimmer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/trimmer"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestTrim_RemovesUnusedKeys(t *testing.T) {
	p := writeTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=secret\nDEBUG=true\n")

	res, err := trimmer.Trim(p, []string{"APP_NAME", "DEBUG"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Removed) != 1 || res.Removed[0] != "DB_PASSWORD" {
		t.Errorf("expected [DB_PASSWORD] removed, got %v", res.Removed)
	}
	if len(res.Kept) != 2 {
		t.Errorf("expected 2 kept entries, got %d", len(res.Kept))
	}
}

func TestTrim_DryRun_DoesNotModifyFile(t *testing.T) {
	original := "APP_NAME=myapp\nOLD_KEY=stale\n"
	p := writeTempEnv(t, original)

	_, err := trimmer.Trim(p, []string{"APP_NAME"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := os.ReadFile(p)
	if string(got) != original {
		t.Errorf("dry-run should not modify file; got %q", string(got))
	}
}

func TestTrim_AllKeysAllowed_NothingRemoved(t *testing.T) {
	p := writeTempEnv(t, "FOO=1\nBAR=2\n")

	res, err := trimmer.Trim(p, []string{"FOO", "BAR"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Removed) != 0 {
		t.Errorf("expected nothing removed, got %v", res.Removed)
	}
}

func TestTrim_NoAllowedKeys_RemovesAll(t *testing.T) {
	p := writeTempEnv(t, "FOO=1\nBAR=2\n")

	res, err := trimmer.Trim(p, []string{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(res.Removed))
	}
	if len(res.Kept) != 0 {
		t.Errorf("expected 0 kept, got %d", len(res.Kept))
	}
}

func TestTrim_InvalidFile_ReturnsError(t *testing.T) {
	_, err := trimmer.Trim("/nonexistent/.env", []string{"FOO"}, false)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
