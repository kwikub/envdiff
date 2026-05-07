package inspector_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/inspector"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestInspect_BasicCounts(t *testing.T) {
	p := writeTempEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_SECRET_KEY=abc\nEMPTY=\n")
	s, err := inspector.Inspect(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalKeys != 4 {
		t.Errorf("TotalKeys: got %d, want 4", s.TotalKeys)
	}
	if s.EmptyValues != 1 {
		t.Errorf("EmptyValues: got %d, want 1", s.EmptyValues)
	}
}

func TestInspect_SecretKeys(t *testing.T) {
	p := writeTempEnv(t, "API_KEY=secret\nDB_PASSWORD=pass\nHOST=localhost\n")
	s, err := inspector.Inspect(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.SecretKeys < 2 {
		t.Errorf("SecretKeys: got %d, want >= 2", s.SecretKeys)
	}
}

func TestInspect_DuplicateKeys(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nFOO=baz\nBAR=qux\n")
	s, err := inspector.Inspect(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.DuplicateKeys) != 1 || s.DuplicateKeys[0] != "FOO" {
		t.Errorf("DuplicateKeys: got %v, want [FOO]", s.DuplicateKeys)
	}
}

func TestInspect_Groups(t *testing.T) {
	p := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=prod\n")
	s, err := inspector.Inspect(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Groups) < 2 {
		t.Errorf("Groups: got %v, want at least DB and APP", s.Groups)
	}
}

func TestInspect_InvalidFile(t *testing.T) {
	_, err := inspector.Inspect("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
