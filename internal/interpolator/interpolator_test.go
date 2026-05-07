package interpolator_test

import (
	"testing"

	"github.com/user/envdiff/internal/interpolator"
	"github.com/user/envdiff/internal/parser"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestInterpolate_NoReferences(t *testing.T) {
	in := entries("HOST", "localhost", "PORT", "5432")
	out, err := interpolator.Interpolate(in, nil, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost" || out[1].Value != "5432" {
		t.Errorf("values should be unchanged, got %v", out)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	in := entries("BASE", "http://example.com", "URL", "${BASE}/api")
	out, err := interpolator.Interpolate(in, nil, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "http://example.com/api" {
		t.Errorf("expected resolved URL, got %q", out[1].Value)
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	in := entries("DOMAIN", "example.com", "EMAIL", "admin@$DOMAIN")
	out, err := interpolator.Interpolate(in, nil, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "admin@example.com" {
		t.Errorf("expected resolved email, got %q", out[1].Value)
	}
}

func TestInterpolate_OverrideWins(t *testing.T) {
	in := entries("URL", "${HOST}/path")
	overrides := map[string]string{"HOST": "override.io"}
	out, err := interpolator.Interpolate(in, overrides, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "override.io/path" {
		t.Errorf("expected override to win, got %q", out[0].Value)
	}
}

func TestInterpolate_StrictMode_ReturnsError(t *testing.T) {
	in := entries("URL", "${MISSING}/path")
	_, err := interpolator.Interpolate(in, nil, interpolator.Options{Strict: true})
	if err == nil {
		t.Fatal("expected error in strict mode for unresolved reference")
	}
}

func TestInterpolate_NonStrictMode_LeavesUnresolved(t *testing.T) {
	in := entries("URL", "${MISSING}/path")
	out, err := interpolator.Interpolate(in, nil, interpolator.Options{Strict: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "${MISSING}/path" {
		t.Errorf("expected original token preserved, got %q", out[0].Value)
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	in := entries("BASE", "http://x.com", "URL", "${BASE}/v1")
	original := in[1].Value
	interpolator.Interpolate(in, nil, interpolator.Options{}) //nolint:errcheck
	if in[1].Value != original {
		t.Error("input entries were mutated")
	}
}
