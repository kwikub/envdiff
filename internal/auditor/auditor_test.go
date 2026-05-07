package auditor_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/auditor"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestAudit_NoIssues(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	result, err := auditor.Audit(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Findings) != 0 {
		t.Errorf("expected no findings, got %d", len(result.Findings))
	}
	if result.HasError {
		t.Error("expected HasError=false")
	}
}

func TestAudit_EmptySecretKey_ReturnsError(t *testing.T) {
	path := writeTempEnv(t, "SECRET_KEY=\nAPP_NAME=myapp\n")
	result, err := auditor.Audit(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasError {
		t.Error("expected HasError=true for empty secret key")
	}
	found := false
	for _, f := range result.Findings {
		if f.Key == "SECRET_KEY" && f.Severity == auditor.SeverityError {
			found = true
		}
	}
	if !found {
		t.Error("expected ERROR finding for SECRET_KEY")
	}
}

func TestAudit_PlaintextSecret_ReturnsInfo(t *testing.T) {
	path := writeTempEnv(t, "API_KEY=abc123\n")
	result, err := auditor.Audit(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range result.Findings {
		if f.Key == "API_KEY" && f.Severity == auditor.SeverityInfo {
			found = true
		}
	}
	if !found {
		t.Error("expected INFO finding for plaintext API_KEY")
	}
}

func TestAudit_UnquotedWhitespace_ReturnsWarning(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=my app\n")
	result, err := auditor.Audit(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range result.Findings {
		if f.Key == "APP_NAME" && f.Severity == auditor.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected WARNING finding for unquoted whitespace in APP_NAME")
	}
}

func TestAudit_InvalidFile_ReturnsError(t *testing.T) {
	_, err := auditor.Audit("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
