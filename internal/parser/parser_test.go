package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(env.Entries))
	}
	if env.Index["DB_HOST"].Value != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", env.Index["DB_HOST"].Value)
	}
}

func TestParse_SkipsCommentsAndBlankLines(t *testing.T) {
	content := "# this is a comment\n\nAPI_KEY=secret\n"
	path := writeTempEnv(t, content)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(env.Entries))
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"\n`)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := env.Index["SECRET"]; ok {
		if v.Value != "my secret value" {
			t.Errorf("expected unquoted value, got %q", v.Value)
		}
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_NO_EQUALS\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
