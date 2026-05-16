package blanker_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/blanker"
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

func TestBlank_AllValues(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	entries, err := blanker.Blank(path, blanker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		if e.Value != "" {
			t.Errorf("key %q: expected blank value, got %q", e.Key, e.Value)
		}
	}
}

func TestBlank_SecretsOnly(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDB_PASSWORD=secret123\nAPI_TOKEN=tok\n")
	entries, err := blanker.Blank(path, blanker.Options{SecretsOnly: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "",
		"API_TOKEN":   "",
	}
	for _, e := range entries {
		if got, ok := expect[e.Key]; ok && got != e.Value {
			t.Errorf("key %q: expected %q, got %q", e.Key, got, e.Value)
		}
	}
}

func TestBlank_ExplicitKeys(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\nDB_PASSWORD=secret\n")
	entries, err := blanker.Blank(path, blanker.Options{Keys: []string{"DB_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range entries {
		switch e.Key {
		case "DB_HOST":
			if e.Value != "" {
				t.Errorf("DB_HOST should be blanked, got %q", e.Value)
			}
		case "APP_NAME", "DB_PASSWORD":
			if e.Value == "" {
				t.Errorf("%q should retain value", e.Key)
			}
		}
	}
}

func TestBlank_Placeholder(t *testing.T) {
	path := writeTempEnv(t, "TOKEN=abc123\n")
	entries, err := blanker.Blank(path, blanker.Options{Placeholder: "REPLACE_ME"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 || entries[0].Value != "REPLACE_ME" {
		t.Errorf("expected placeholder value, got %q", entries[0].Value)
	}
}

func TestBlank_InvalidFile(t *testing.T) {
	_, err := blanker.Blank("/nonexistent/path/.env", blanker.Options{})
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestFormat_ProducesEnvLines(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\nSECRET=hidden\n")
	entries, _ := blanker.Blank(path, blanker.Options{})
	out := blanker.Format(entries)
	if out == "" {
		t.Error("expected non-empty formatted output")
	}
	for _, e := range entries {
		expected := e.Key + "=\n"
		if !containsLine(out, expected) {
			t.Errorf("formatted output missing line %q", expected)
		}
	}
}

func containsLine(s, line string) bool {
	for _, l := range splitLines(s) {
		if l == line {
			return true
		}
	}
	return false
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i+1])
			start = i + 1
		}
	}
	return lines
}
