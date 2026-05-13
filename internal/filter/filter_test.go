package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func keys(es []parser.Entry) []string {
	var ks []string
	for _, e := range es {
		ks = append(ks, e.Key)
	}
	return ks
}

func TestFilter_ByPrefix(t *testing.T) {
	input := entries("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "myapp")
	got := filter.Filter(input, filter.Options{Prefix: "DB_"})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key != "DB_HOST" && e.Key != "DB_PORT" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestFilter_ByAllowList(t *testing.T) {
	input := entries("A", "1", "B", "2", "C", "3")
	got := filter.Filter(input, filter.Options{Keys: []string{"A", "C"}})
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilter_ByExclude(t *testing.T) {
	input := entries("A", "1", "B", "2", "C", "3")
	got := filter.Filter(input, filter.Options{Exclude: []string{"B"}})
	for _, e := range got {
		if e.Key == "B" {
			t.Error("excluded key B should not be present")
		}
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilter_SecretsOnly(t *testing.T) {
	input := entries("APP_NAME", "myapp", "DB_PASSWORD", "s3cr3t", "API_TOKEN", "tok", "PORT", "8080")
	got := filter.Filter(input, filter.Options{SecretsOnly: true})
	k := keys(got)
	if len(k) != 2 {
		t.Fatalf("expected 2 secret entries, got %d: %v", len(k), k)
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	input := entries("X", "1", "Y", "2")
	got := filter.Filter(input, filter.Options{})
	if len(got) != len(input) {
		t.Errorf("expected all %d entries, got %d", len(input), len(got))
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	got := filter.Filter(nil, filter.Options{Prefix: "DB_"})
	if got != nil && len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
