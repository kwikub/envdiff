package rotator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/rotator"
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

func TestRotate_SecretKeysChanged(t *testing.T) {
	p := writeTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=old_pass\nAPI_TOKEN=old_token\n")

	results, err := rotator.Rotate(p, rotator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 rotations, got %d", len(results))
	}
	for _, r := range results {
		if r.NewValue == r.OldValue {
			t.Errorf("key %q: new value equals old value", r.Key)
		}
		if len(r.NewValue) != 32 { // default 16 bytes → 32 hex chars
			t.Errorf("key %q: expected 32-char hex, got %q", r.Key, r.NewValue)
		}
	}
}

func TestRotate_NonSecretKeysUnchanged(t *testing.T) {
	p := writeTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=old_pass\n")

	_, err := rotator.Rotate(p, rotator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, _ := parser.Parse(p)
	for _, e := range entries {
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("APP_NAME should not have been rotated, got %q", e.Value)
		}
	}
}

func TestRotate_ExplicitKeys(t *testing.T) {
	p := writeTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=old_pass\nAPI_TOKEN=old_token\n")

	results, err := rotator.Rotate(p, rotator.Options{Keys: []string{"API_TOKEN"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "API_TOKEN" {
		t.Fatalf("expected only API_TOKEN to be rotated, got %+v", results)
	}
}

func TestRotate_DryRun_DoesNotModifyFile(t *testing.T) {
	original := "DB_PASSWORD=original\n"
	p := writeTempEnv(t, original)

	_, err := rotator.Rotate(p, rotator.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(p)
	if string(data) != original {
		t.Errorf("dry-run modified file: got %q", string(data))
	}
}

func TestRotate_CustomLength(t *testing.T) {
	p := writeTempEnv(t, "API_KEY=old\n")

	results, err := rotator.Rotate(p, rotator.Options{Length: 8})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if len(results[0].NewValue) != 16 { // 8 bytes → 16 hex chars
		t.Errorf("expected 16-char hex for length=8, got %q", results[0].NewValue)
	}
}
