package scoper_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/scoper"
)

type Entry = scoper.Entry

func entries(kvs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestScope_NoScopeOption_ReturnsAll(t *testing.T) {
	in := entries("A", "1", "B", "2")
	out, err := scoper.Scope(in, scoper.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestScope_MatchingTag_Included(t *testing.T) {
	in := []Entry{
		{Key: "DB_URL", Value: "x", Comment: "@scope:production"},
		{Key: "API_KEY", Value: "y", Comment: "@scope:staging"},
		{Key: "LOG_LEVEL", Value: "info"}, // no tag
	}
	out, err := scoper.Scope(in, scoper.Options{Scope: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries (matching + untagged), got %d", len(out))
	}
	if out[0].Key != "DB_URL" {
		t.Errorf("expected DB_URL, got %s", out[0].Key)
	}
	if out[1].Key != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL, got %s", out[1].Key)
	}
}

func TestScope_StrictMode_WrongScope_ReturnsError(t *testing.T) {
	in := []Entry{
		{Key: "SECRET", Value: "val", Comment: "@scope:staging"},
	}
	_, err := scoper.Scope(in, scoper.Options{Scope: "production", Strict: true})
	if err == nil {
		t.Fatal("expected error for mismatched scope in strict mode")
	}
}

func TestScope_CustomTagPrefix(t *testing.T) {
	in := []Entry{
		{Key: "X", Value: "1", Comment: "env:production"},
		{Key: "Y", Value: "2", Comment: "env:staging"},
	}
	out, err := scoper.Scope(in, scoper.Options{Scope: "production", TagPrefix: "env:"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Key != "X" {
		t.Errorf("expected only X, got %+v", out)
	}
}

func TestTag_AppendsScope(t *testing.T) {
	in := entries("DB_URL", "x", "LOG_LEVEL", "info")
	out := scoper.Tag(in, "production", []string{"DB_URL"}, "")
	if len(out) != 2 {
		t.Fatalf("expected 2 entries")
	}
	if out[0].Comment == "" {
		t.Errorf("expected scope tag on DB_URL comment, got empty")
	}
	if out[1].Comment != "" {
		t.Errorf("expected no tag on LOG_LEVEL, got %q", out[1].Comment)
	}
}

func TestTag_ReplacesExistingScope(t *testing.T) {
	in := []Entry{{Key: "DB_URL", Value: "x", Comment: "@scope:staging"}}
	out := scoper.Tag(in, "production", []string{"DB_URL"}, "")
	scope, ok := extractScope(out[0].Comment, "@scope:")
	if !ok || scope != "production" {
		t.Errorf("expected scope=production, got comment=%q", out[0].Comment)
	}
}

func extractScope(comment, prefix string) (string, bool) {
	import_idx := -1
	for i := 0; i <= len(comment)-len(prefix); i++ {
		if comment[i:i+len(prefix)] == prefix {
			import_idx = i
			break
		}
	}
	if import_idx == -1 {
		return "", false
	}
	rest := comment[import_idx+len(prefix):]
	for _, f := range splitFields(rest) {
		return f, true
	}
	return "", false
}

func splitFields(s string) []string {
	var out []string
	for _, f := range []string{} {
		_ = f
	}
	// simple split
	word := ""
	for _, c := range s {
		if c == ' ' || c == '\t' {
			if word != "" {
				out = append(out, word)
				word = ""
			}
		} else {
			word += string(c)
		}
	}
	if word != "" {
		out = append(out, word)
	}
	return out
}
