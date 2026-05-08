package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/redactor"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestRedact_SecretValuesReplaced(t *testing.T) {
	in := entries("API_KEY", "abc123", "APP_NAME", "myapp")
	out := redactor.Redact(in, redactor.Options{})

	if out[0].Value != "[REDACTED]" {
		t.Errorf("expected API_KEY to be redacted, got %q", out[0].Value)
	}
	if out[1].Value != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", out[1].Value)
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	in := entries("SECRET", "s3cr3t")
	origVal := in[0].Value
	redactor.Redact(in, redactor.Options{})
	if in[0].Value != origVal {
		t.Error("original slice was mutated")
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	in := entries("PASSWORD", "hunter2")
	out := redactor.Redact(in, redactor.Options{Placeholder: "***"})
	if out[0].Value != "***" {
		t.Errorf("expected custom placeholder, got %q", out[0].Value)
	}
}

func TestRedact_ExtraKeys(t *testing.T) {
	in := entries("MY_TOKEN", "tok", "REGION", "us-east-1")
	out := redactor.Redact(in, redactor.Options{ExtraKeys: []string{"REGION"}})
	if out[0].Value != "[REDACTED]" {
		t.Errorf("MY_TOKEN should be redacted")
	}
	if out[1].Value != "[REDACTED]" {
		t.Errorf("REGION should be redacted via ExtraKeys")
	}
}

func TestRedactMap_SecretValuesReplaced(t *testing.T) {
	m := map[string]string{"DB_PASSWORD": "secret", "HOST": "localhost"}
	out := redactor.RedactMap(m, redactor.Options{})
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", out["DB_PASSWORD"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST unchanged, got %q", out["HOST"])
	}
}

func TestRedactMap_DoesNotMutateOriginal(t *testing.T) {
	m := map[string]string{"API_SECRET": "original"}
	redactor.RedactMap(m, redactor.Options{})
	if m["API_SECRET"] != "original" {
		t.Error("original map was mutated")
	}
}
