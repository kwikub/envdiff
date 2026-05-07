package template_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/template"
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

func TestGenerate_BlankValues(t *testing.T) {
	src := writeTempEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\n")
	out := filepath.Join(t.TempDir(), "out.env")

	if err := template.Generate(src, out, template.Options{}); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, _ := os.ReadFile(out)
	got := string(data)
	if !strings.Contains(got, "APP_HOST=") {
		t.Errorf("expected APP_HOST= in output, got:\n%s", got)
	}
	if strings.Contains(got, "localhost") {
		t.Errorf("expected value to be stripped, got:\n%s", got)
	}
}

func TestGenerate_MaskSecrets(t *testing.T) {
	src := writeTempEnv(t, "API_KEY=supersecret\nAPP_NAME=myapp\n")
	out := filepath.Join(t.TempDir(), "out.env")

	if err := template.Generate(src, out, template.Options{MaskSecrets: true}); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, _ := os.ReadFile(out)
	got := string(data)
	if !strings.Contains(got, "API_KEY=<secret>") {
		t.Errorf("expected API_KEY=<secret>, got:\n%s", got)
	}
	if strings.Contains(got, "supersecret") {
		t.Errorf("secret value should not appear in output, got:\n%s", got)
	}
}

func TestGenerate_AddComments(t *testing.T) {
	src := writeTempEnv(t, "DB_URL=postgres://localhost/db\nSECRET_KEY=abc123\n")
	out := filepath.Join(t.TempDir(), "out.env")

	if err := template.Generate(src, out, template.Options{AddComments: true}); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, _ := os.ReadFile(out)
	got := string(data)
	if !strings.Contains(got, "# DB_URL") {
		t.Errorf("expected comment for DB_URL, got:\n%s", got)
	}
	if !strings.Contains(got, "# SECRET_KEY") {
		t.Errorf("expected comment for SECRET_KEY, got:\n%s", got)
	}
}

func TestGenerate_InvalidSource(t *testing.T) {
	err := template.Generate("/nonexistent/path.env", "/tmp/out.env", template.Options{})
	if err == nil {
		t.Error("expected error for missing source file, got nil")
	}
}
