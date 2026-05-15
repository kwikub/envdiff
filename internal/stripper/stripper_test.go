package stripper_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/stripper"
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

func TestStrip_RemovesCommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\nFOO=bar\n\nBAZ=qux\n")
	lines, err := stripper.Strip(path, stripper.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "FOO=bar" || lines[1] != "BAZ=qux" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestStrip_KeepsCommentsWhenDisabled(t *testing.T) {
	path := writeTempEnv(t, "# a comment\nFOO=bar\n")
	opts := stripper.Options{RemoveComments: false, RemoveBlankLines: true}
	lines, err := stripper.Strip(path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestStrip_KeepsBlanksWhenDisabled(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\n\nBAZ=qux\n")
	opts := stripper.Options{RemoveComments: true, RemoveBlankLines: false}
	lines, err := stripper.Strip(path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (including blank), got %d", len(lines))
	}
}

func TestStrip_RemovesDisabledEntries(t *testing.T) {
	path := writeTempEnv(t, "# SECRET=old\nFOO=bar\n# just a comment\n")
	opts := stripper.Options{
		RemoveComments:   false,
		RemoveBlankLines: true,
		RemoveDisabled:   true,
	}
	lines, err := stripper.Strip(path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "# SECRET=old" removed (disabled entry), "FOO=bar" kept, "# just a comment" kept
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
}

func TestStrip_InvalidPath_ReturnsError(t *testing.T) {
	_, err := stripper.Strip("/nonexistent/path/.env", stripper.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
