package normalizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/normalizer"
	"github.com/user/envdiff/internal/parser"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	in := entries("db_host", "localhost", "api_key", "secret")
	out := normalizer.Normalize(in, normalizer.Options{UppercaseKeys: true})

	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out[0].Key)
	}
	if out[1].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", out[1].Key)
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	in := entries("HOST", "  localhost  ", "PORT", "\t8080\t")
	out := normalizer.Normalize(in, normalizer.Options{TrimValues: true})

	if out[0].Value != "localhost" {
		t.Errorf("expected 'localhost', got %q", out[0].Value)
	}
	if out[1].Value != "8080" {
		t.Errorf("expected '8080', got %q", out[1].Value)
	}
}

func TestNormalize_StripQuotes_Double(t *testing.T) {
	in := entries("NAME", `"alice"`)
	out := normalizer.Normalize(in, normalizer.Options{StripQuotes: true})

	if out[0].Value != "alice" {
		t.Errorf("expected alice, got %q", out[0].Value)
	}
}

func TestNormalize_StripQuotes_Single(t *testing.T) {
	in := entries("NAME", "'bob'")
	out := normalizer.Normalize(in, normalizer.Options{StripQuotes: true})

	if out[0].Value != "bob" {
		t.Errorf("expected bob, got %q", out[0].Value)
	}
}

func TestNormalize_StripQuotes_MismatchedUnchanged(t *testing.T) {
	in := entries("NAME", `"mixed'`)
	out := normalizer.Normalize(in, normalizer.Options{StripQuotes: true})

	if out[0].Value != `"mixed'` {
		t.Errorf("expected unchanged value, got %q", out[0].Value)
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	in := entries("key", "  value  ")
	orig := in[0].Value
	normalizer.Normalize(in, normalizer.Options{TrimValues: true, UppercaseKeys: true})

	if in[0].Value != orig {
		t.Errorf("original entry was mutated: got %q", in[0].Value)
	}
	if in[0].Key != "key" {
		t.Errorf("original key was mutated: got %s", in[0].Key)
	}
}

func TestNormalize_AllOptionsOff_ReturnsUnchanged(t *testing.T) {
	in := entries("my_key", "  'value'  ")
	out := normalizer.Normalize(in, normalizer.Options{})

	if out[0].Key != "my_key" || out[0].Value != "  'value'  " {
		t.Errorf("unexpected change: %+v", out[0])
	}
}
