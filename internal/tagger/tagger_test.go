package tagger_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/tagger"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestApply_NoOptions_ReturnsResultsWithNoTags(t *testing.T) {
	e := entries("APP_NAME", "myapp", "PORT", "8080")
	res, err := tagger.Apply(e, tagger.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}
	for _, r := range res {
		if len(r.Tags) != 0 {
			t.Errorf("key %s: expected no tags, got %v", r.Key, r.Tags)
		}
	}
}

func TestApply_TagSecrets_MarksSecretKeys(t *testing.T) {
	e := entries("DB_PASSWORD", "s3cr3t", "APP_NAME", "myapp", "API_TOKEN", "tok")
	res, err := tagger.Apply(e, tagger.Options{TagSecrets: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tagged := map[string][]string{}
	for _, r := range res {
		tagged[r.Key] = r.Tags
	}
	if !containsTag(tagged["DB_PASSWORD"], "secret") {
		t.Error("DB_PASSWORD should be tagged secret")
	}
	if !containsTag(tagged["API_TOKEN"], "secret") {
		t.Error("API_TOKEN should be tagged secret")
	}
	if containsTag(tagged["APP_NAME"], "secret") {
		t.Error("APP_NAME should not be tagged secret")
	}
}

func TestApply_TagEmpty_MarksBlankValues(t *testing.T) {
	e := entries("FILLED", "yes", "EMPTY_VAR", "", "SPACES", "   ")
	res, err := tagger.Apply(e, tagger.Options{TagEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tagged := map[string][]string{}
	for _, r := range res {
		tagged[r.Key] = r.Tags
	}
	if containsTag(tagged["FILLED"], "empty") {
		t.Error("FILLED should not be tagged empty")
	}
	if !containsTag(tagged["EMPTY_VAR"], "empty") {
		t.Error("EMPTY_VAR should be tagged empty")
	}
	if !containsTag(tagged["SPACES"], "empty") {
		t.Error("SPACES should be tagged empty")
	}
}

func TestApply_ExtraTags_AppliedToSpecifiedKeys(t *testing.T) {
	e := entries("SERVICE_A", "a", "SERVICE_B", "b", "UNRELATED", "x")
	opts := tagger.Options{
		Extra: map[string][]string{
			"critical": {"SERVICE_A", "SERVICE_B"},
		},
	}
	res, err := tagger.Apply(e, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range res {
		switch r.Key {
		case "SERVICE_A", "SERVICE_B":
			if !containsTag(r.Tags, "critical") {
				t.Errorf("%s should be tagged critical", r.Key)
			}
		case "UNRELATED":
			if containsTag(r.Tags, "critical") {
				t.Error("UNRELATED should not be tagged critical")
			}
		}
	}
}

func TestApply_NilEntries_ReturnsError(t *testing.T) {
	_, err := tagger.Apply(nil, tagger.Options{})
	if err == nil {
		t.Fatal("expected error for nil entries")
	}
}

func TestFilterByTag_ReturnsOnlyMatchingResults(t *testing.T) {
	e := entries("DB_PASSWORD", "x", "PORT", "3000", "API_SECRET", "y")
	res, _ := tagger.Apply(e, tagger.Options{TagSecrets: true})
	secrets := tagger.FilterByTag(res, "secret")
	if len(secrets) != 2 {
		t.Fatalf("expected 2 secret results, got %d", len(secrets))
	}
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
