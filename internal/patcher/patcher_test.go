package patcher_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/patcher"
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

func TestApply_SetNewKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := patcher.Apply(path, []patcher.Patch{
		{Op: patcher.OpSet, Key: "NEW_KEY", Value: "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 {
		t.Fatalf("expected 1 applied, got %d", len(res.Applied))
	}
	out := patcher.Format(res.Entries)
	if !strings.Contains(out, "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY=hello in output, got:\n%s", out)
	}
}

func TestApply_UpdateExistingKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=old\nBAR=keep\n")
	res, err := patcher.Apply(path, []patcher.Patch{
		{Op: patcher.OpSet, Key: "FOO", Value: "new"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := patcher.Format(res.Entries)
	if !strings.Contains(out, "FOO=new") {
		t.Errorf("expected FOO=new, got:\n%s", out)
	}
	if strings.Contains(out, "FOO=old") {
		t.Errorf("old value should be gone, got:\n%s", out)
	}
}

func TestApply_DeleteExistingKey(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nSECRET=topsecret\n")
	res, err := patcher.Apply(path, []patcher.Patch{
		{Op: patcher.OpDelete, Key: "SECRET"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Applied) != 1 {
		t.Fatalf("expected 1 applied, got %d", len(res.Applied))
	}
	out := patcher.Format(res.Entries)
	if strings.Contains(out, "SECRET") {
		t.Errorf("SECRET should be removed, got:\n%s", out)
	}
}

func TestApply_DeleteMissingKey_Skipped(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	res, err := patcher.Apply(path, []patcher.Patch{
		{Op: patcher.OpDelete, Key: "GHOST"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if len(res.Applied) != 0 {
		t.Fatalf("expected 0 applied, got %d", len(res.Applied))
	}
}

func TestApply_UnknownOp_ReturnsError(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n")
	_, err := patcher.Apply(path, []patcher.Patch{
		{Op: "upsert", Key: "FOO", Value: "val"},
	})
	if err == nil {
		t.Fatal("expected error for unknown op, got nil")
	}
}

func TestFormat_ProducesKeyEqValue(t *testing.T) {
	path := writeTempEnv(t, "A=1\nB=2\n")
	res, err := patcher.Apply(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := patcher.Format(res.Entries)
	for _, line := range []string{"A=1", "B=2"} {
		if !strings.Contains(out, line) {
			t.Errorf("expected %q in output:\n%s", line, out)
		}
	}
}
