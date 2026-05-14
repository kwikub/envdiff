package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/resolver"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	return filepath.Clean(f.Name())
}

func TestResolve_NoFiles_ReturnsError(t *testing.T) {
	_, err := resolver.Resolve(nil, resolver.StrategyLast)
	if err == nil {
		t.Fatal("expected error for empty file list")
	}
}

func TestResolve_UnknownStrategy_ReturnsError(t *testing.T) {
	f := writeTempEnv(t, "KEY=val\n")
	_, err := resolver.Resolve([]string{f}, "unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestResolve_SingleFile_ReturnsAllKeys(t *testing.T) {
	f := writeTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	results, err := resolver.Resolve([]string{f}, resolver.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestResolve_StrategyFirst_KeepsFirstValue(t *testing.T) {
	f1 := writeTempEnv(t, "DB_HOST=primary\n")
	f2 := writeTempEnv(t, "DB_HOST=replica\n")

	results, err := resolver.Resolve([]string{f1, f2}, resolver.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != "primary" {
		t.Errorf("expected 'primary', got %q", results[0].Value)
	}
	if !results[0].Overridden {
		t.Error("expected Overridden=true")
	}
}

func TestResolve_StrategyLast_KeepsLastValue(t *testing.T) {
	f1 := writeTempEnv(t, "DB_HOST=primary\n")
	f2 := writeTempEnv(t, "DB_HOST=replica\n")

	results, err := resolver.Resolve([]string{f1, f2}, resolver.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Value != "replica" {
		t.Errorf("expected 'replica', got %q", results[0].Value)
	}
	if results[0].Source != f2 {
		t.Errorf("expected source %q, got %q", f2, results[0].Source)
	}
}

func TestResolve_UniqueKeys_NotMarkedOverridden(t *testing.T) {
	f1 := writeTempEnv(t, "KEY_A=foo\n")
	f2 := writeTempEnv(t, "KEY_B=bar\n")

	results, err := resolver.Resolve([]string{f1, f2}, resolver.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Overridden {
			t.Errorf("key %q should not be marked overridden", r.Key)
		}
	}
}
