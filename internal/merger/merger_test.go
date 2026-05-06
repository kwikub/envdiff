package merger_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/merger"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestMerge_NoFiles(t *testing.T) {
	_, err := merger.Merge(nil, merger.StrategyLast)
	if err == nil {
		t.Fatal("expected error for empty file list")
	}
}

func TestMerge_SingleFile(t *testing.T) {
	f := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	res, err := merger.Merge([]string{f}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Key != "FOO" || res.Entries[0].Value != "bar" {
		t.Errorf("unexpected entry: %+v", res.Entries[0])
	}
}

func TestMerge_StrategyLast_WinsOnConflict(t *testing.T) {
	f1 := writeTempEnv(t, "FOO=first\nSHARED=one\n")
	f2 := writeTempEnv(t, "SHARED=two\nBAR=second\n")
	res, err := merger.Merge([]string{f1, f2}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range res.Entries {
		if e.Key == "SHARED" && e.Value != "two" {
			t.Errorf("StrategyLast: expected SHARED=two, got %q", e.Value)
		}
	}
}

func TestMerge_StrategyFirst_WinsOnConflict(t *testing.T) {
	f1 := writeTempEnv(t, "SHARED=one\n")
	f2 := writeTempEnv(t, "SHARED=two\n")
	res, err := merger.Merge([]string{f1, f2}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range res.Entries {
		if e.Key == "SHARED" && e.Value != "one" {
			t.Errorf("StrategyFirst: expected SHARED=one, got %q", e.Value)
		}
	}
}

func TestMerge_RecordsConflicts(t *testing.T) {
	f1 := writeTempEnv(t, "KEY=alpha\n")
	f2 := writeTempEnv(t, "KEY=beta\n")
	res, err := merger.Merge([]string{f1, f2}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Conflicts[0].Key != "KEY" {
		t.Errorf("expected conflict key KEY, got %q", res.Conflicts[0].Key)
	}
}

func TestMerge_UniqueKeysAcrossFiles(t *testing.T) {
	f1 := writeTempEnv(t, "A=1\nB=2\n")
	f2 := writeTempEnv(t, "C=3\nD=4\n")
	res, err := merger.Merge([]string{f1, f2}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(res.Entries))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}
