package scanner_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/scanner"
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

func TestScan_NoFindings(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	findings, err := scanner.Scan(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestScan_HighEntropyToken(t *testing.T) {
	path := writeTempEnv(t, "API_KEY=aB3dEfGhIjKlMnOpQrStUvWxYz1234567890abcd\n")
	findings, err := scanner.Scan(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Rule == "high-entropy" {
			found = true
			if f.Severity != scanner.SeverityWarning {
				t.Errorf("expected WARNING severity, got %s", f.Severity)
			}
			if f.Value != "***" {
				t.Errorf("expected masked value, got %q", f.Value)
			}
		}
	}
	if !found {
		t.Error("expected high-entropy finding")
	}
}

func TestScan_URLWithCredentials(t *testing.T) {
	path := writeTempEnv(t, "DATABASE_URL=postgres://admin:s3cr3t@localhost/db\n")
	findings, err := scanner.Scan(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Rule == "url-credentials" {
			found = true
			if f.Severity != scanner.SeverityError {
				t.Errorf("expected ERROR severity, got %s", f.Severity)
			}
		}
	}
	if !found {
		t.Error("expected url-credentials finding")
	}
}

func TestScan_PlaceholderValue(t *testing.T) {
	path := writeTempEnv(t, "DB_PASSWORD=changeme\n")
	findings, err := scanner.Scan(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, f := range findings {
		if f.Rule == "placeholder-value" {
			found = true
			if f.Severity != scanner.SeverityInfo {
				t.Errorf("expected INFO severity, got %s", f.Severity)
			}
		}
	}
	if !found {
		t.Error("expected placeholder-value finding")
	}
}

func TestScan_InvalidFile_ReturnsError(t *testing.T) {
	_, err := scanner.Scan("/nonexistent/path.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
