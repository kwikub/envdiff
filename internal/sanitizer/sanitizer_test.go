package sanitizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/sanitizer"
)

func entries(kv ...string) []parser.Entry {
	out := make([]parser.Entry, 0, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		out = append(out, parser.Entry{Key: kv[i], Value: kv[i+1]})
	}
	return out
}

func TestSanitize_TrimValues(t *testing.T) {
	in := entries("HOST", "  localhost  ", "PORT", "\t8080\t")
	res, err := sanitizer.Sanitize(in, sanitizer.Options{TrimValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Entry.Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", res[0].Entry.Value)
	}
	if res[1].Entry.Value != "8080" {
		t.Errorf("expected '8080', got %q", res[1].Entry.Value)
	}
}

func TestSanitize_UppercaseKeys(t *testing.T) {
	in := entries("db_host", "localhost", "api_key", "secret")
	res, err := sanitizer.Sanitize(in, sanitizer.Options{UppercaseKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Entry.Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", res[0].Entry.Key)
	}
	if res[1].Entry.Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %q", res[1].Entry.Key)
	}
}

func TestSanitize_StripQuotes_DoubleQuoted(t *testing.T) {
	in := entries("MSG", `"hello world"`)
	res, err := sanitizer.Sanitize(in, sanitizer.Options{StripQuotes: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Entry.Value != "hello world" {
		t.Errorf("expected 'hello world', got %q", res[0].Entry.Value)
	}
}

func TestSanitize_StripQuotes_SingleQuoted(t *testing.T) {
	in := entries("MSG", "'hello'")
	res, err := sanitizer.Sanitize(in, sanitizer.Options{StripQuotes: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Entry.Value != "hello" {
		t.Errorf("expected 'hello', got %q", res[0].Entry.Value)
	}
}

func TestSanitize_RejectInvalid_ReturnsError(t *testing.T) {
	in := entries("123BAD", "value")
	_, err := sanitizer.Sanitize(in, sanitizer.Options{RejectInvalid: true})
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

func TestSanitize_InvalidKey_Warning_NoReject(t *testing.T) {
	in := entries("123BAD", "value")
	res, err := sanitizer.Sanitize(in, sanitizer.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Warning == "" {
		t.Error("expected a warning for invalid key name, got empty string")
	}
}

func TestSanitize_ValidKey_NoWarning(t *testing.T) {
	in := entries("VALID_KEY", "value")
	res, err := sanitizer.Sanitize(in, sanitizer.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res[0].Warning != "" {
		t.Errorf("expected no warning, got %q", res[0].Warning)
	}
}

func TestSanitize_PreservesOrder(t *testing.T) {
	in := entries("Z_KEY", "z", "A_KEY", "a", "M_KEY", "m")
	res, err := sanitizer.Sanitize(in, sanitizer.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := []string{res[0].Entry.Key, res[1].Entry.Key, res[2].Entry.Key}
	expected := []string{"Z_KEY", "A_KEY", "M_KEY"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
