package transformer_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/transformer"
)

func entries(kv ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kv); i += 2 {
		out = append(out, parser.Entry{Key: kv[i], Value: kv[i+1]})
	}
	return out
}

func TestApply_UppercaseKeys(t *testing.T) {
	in := entries("db_host", "localhost", "db_port", "5432")
	out, err := transformer.Apply(in, transformer.Options{
		Transforms: []transformer.Transform{transformer.UppercaseKeys},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "DB_HOST" || out[1].Key != "DB_PORT" {
		t.Errorf("expected uppercase keys, got %v", out)
	}
}

func TestApply_LowercaseKeys(t *testing.T) {
	in := entries("DB_HOST", "localhost")
	out, err := transformer.Apply(in, transformer.Options{
		Transforms: []transformer.Transform{transformer.LowercaseKeys},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "db_host" {
		t.Errorf("expected lowercase key, got %q", out[0].Key)
	}
}

func TestApply_TrimValues(t *testing.T) {
	in := entries("HOST", "  localhost  ")
	out, err := transformer.Apply(in, transformer.Options{
		Transforms: []transformer.Transform{transformer.TrimValues},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost" {
		t.Errorf("expected trimmed value, got %q", out[0].Value)
	}
}

func TestApply_AddPrefix(t *testing.T) {
	in := entries("HOST", "localhost", "APP_PORT", "8080")
	out, err := transformer.Apply(in, transformer.Options{
		Transforms: []transformer.Transform{transformer.AddPrefix},
		Prefix:     "APP_",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %q", out[0].Key)
	}
	// already prefixed — should not double-prefix
	if out[1].Key != "APP_PORT" {
		t.Errorf("expected APP_PORT unchanged, got %q", out[1].Key)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	in := entries("APP_HOST", "localhost", "HOST", "other")
	out, err := transformer.Apply(in, transformer.Options{
		Transforms: []transformer.Transform{transformer.StripPrefix},
		Prefix:     "APP_",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "HOST" {
		t.Errorf("expected HOST, got %q", out[0].Key)
	}
	if out[1].Key != "HOST" {
		t.Errorf("expected HOST unchanged, got %q", out[1].Key)
	}
}

func TestApply_UnknownTransform_ReturnsError(t *testing.T) {
	in := entries("KEY", "val")
	_, err := transformer.Apply(in, transformer.Options{
		Transforms: []transformer.Transform{"bogus"},
	})
	if err == nil {
		t.Error("expected error for unknown transform")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	in := entries("key", "value")
	orig := in[0].Key
	_, _ = transformer.Apply(in, transformer.Options{
		Transforms: []transformer.Transform{transformer.UppercaseKeys},
	})
	if in[0].Key != orig {
		t.Error("original entries were mutated")
	}
}
