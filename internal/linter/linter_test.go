package linter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestLint_NoIssues(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	findings, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=first\nAPP_NAME=second\n")
	findings, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != linter.SeverityError {
		t.Errorf("expected error severity, got %s", findings[0].Severity)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=\n")
	findings, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestLint_InvalidKeyName(t *testing.T) {
	path := writeTempEnv(t, "app_name=myapp\n")
	findings, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "app_name" {
		t.Errorf("unexpected key in finding: %s", findings[0].Key)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	content := "app_name=\napp_name=\n"
	path := writeTempEnv(t, content)
	findings, err := linter.Lint(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Each line: empty value (warning) + bad name (warning); second line also duplicate (error)
	if len(findings) < 3 {
		t.Errorf("expected at least 3 findings, got %d: %v", len(findings), findings)
	}
}
