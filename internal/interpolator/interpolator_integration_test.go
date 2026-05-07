package interpolator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/interpolator"
	"github.com/user/envdiff/internal/parser"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestInterpolate_FromFile_ChainedRefs(t *testing.T) {
	path := writeTempEnv(t, "PROTO=https\nHOST=api.example.com\nBASE_URL=${PROTO}://${HOST}\nFULL_URL=${BASE_URL}/v2\n")
	entries, err := parser.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := interpolator.Interpolate(entries, nil, interpolator.Options{})
	if err != nil {
		t.Fatalf("interpolate: %v", err)
	}
	result := toMap(out)
	if result["BASE_URL"] != "https://api.example.com" {
		t.Errorf("BASE_URL: got %q", result["BASE_URL"])
	}
	// FULL_URL references BASE_URL which itself contained references;
	// single-pass resolution means it resolves the already-resolved BASE_URL.
	if result["FULL_URL"] != "https://api.example.com/v2" {
		t.Errorf("FULL_URL: got %q", result["FULL_URL"])
	}
}

func TestInterpolate_FallbackToOS(t *testing.T) {
	t.Setenv("OS_HOST", "os-host.internal")
	path := writeTempEnv(t, "URL=http://${OS_HOST}/path\n")
	entries, err := parser.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := interpolator.Interpolate(entries, nil, interpolator.Options{FallbackToOS: true})
	if err != nil {
		t.Fatalf("interpolate: %v", err)
	}
	if out[0].Value != "http://os-host.internal/path" {
		t.Errorf("expected OS fallback, got %q", out[0].Value)
	}
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
