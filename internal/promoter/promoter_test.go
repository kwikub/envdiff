package promoter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/promoter"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestPromote_BasicPromotion(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := writeTempEnv(t, "FOO=old\n")

	results, err := promoter.Promote(src, dst, promoter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Action != "promoted" {
			t.Errorf("key %q: expected promoted, got %q", r.Key, r.Action)
		}
	}
}

func TestPromote_SkipExisting(t *testing.T) {
	src := writeTempEnv(t, "FOO=newval\nBAR=baz\n")
	dst := writeTempEnv(t, "FOO=original\n")

	results, err := promoter.Promote(src, dst, promoter.Options{SkipExisting: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var skipped, promoted int
	for _, r := range results {
		switch r.Action {
		case "skipped":
			skipped++
		case "promoted":
			promoted++
		}
	}
	if skipped != 1 || promoted != 1 {
		t.Errorf("expected 1 skipped and 1 promoted, got %d skipped %d promoted", skipped, promoted)
	}
}

func TestPromote_AllowList(t *testing.T) {
	src := writeTempEnv(t, "FOO=1\nBAR=2\nBAZ=3\n")
	dst := writeTempEnv(t, "")

	results, err := promoter.Promote(src, dst, promoter.Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Key == "BAR" {
			t.Error("BAR should not have been promoted")
		}
	}
}

func TestPromote_DryRun_DoesNotModifyFile(t *testing.T) {
	src := writeTempEnv(t, "FOO=new\n")
	dst := writeTempEnv(t, "FOO=original\n")

	original, _ := os.ReadFile(dst)
	results, err := promoter.Promote(src, dst, promoter.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after, _ := os.ReadFile(dst)
	if string(original) != string(after) {
		t.Error("dry-run should not modify the destination file")
	}
	if len(results) == 0 || results[0].Action != "dry-run" {
		t.Error("expected dry-run action in results")
	}
}

func TestPromote_InvalidSrc_ReturnsError(t *testing.T) {
	dst := writeTempEnv(t, "FOO=bar\n")
	_, err := promoter.Promote(filepath.Join(t.TempDir(), "missing.env"), dst, promoter.Options{})
	if err == nil {
		t.Error("expected error for missing src file")
	}
}
