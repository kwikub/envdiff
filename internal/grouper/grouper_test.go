package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/grouper"
	"github.com/user/envdiff/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	out := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestGroupBy_BasicPrefixes(t *testing.T) {
	in := entries("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "envdiff")
	groups := grouper.GroupBy(in, grouper.Options{})

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "APP" {
		t.Errorf("expected first group APP, got %s", groups[0].Name)
	}
	if groups[1].Name != "DB" {
		t.Errorf("expected second group DB, got %s", groups[1].Name)
	}
	if len(groups[1].Entries) != 2 {
		t.Errorf("expected 2 entries in DB group, got %d", len(groups[1].Entries))
	}
}

func TestGroupBy_UngroupedIncluded(t *testing.T) {
	in := entries("DB_HOST", "localhost", "PORT", "8080")
	groups := grouper.GroupBy(in, grouper.Options{IncludeUngrouped: true})

	var found bool
	for _, g := range groups {
		if g.Name == "(ungrouped)" {
			found = true
			if len(g.Entries) != 1 {
				t.Errorf("expected 1 ungrouped entry, got %d", len(g.Entries))
			}
		}
	}
	if !found {
		t.Error("expected (ungrouped) group to be present")
	}
}

func TestGroupBy_UngroupedExcluded(t *testing.T) {
	in := entries("DB_HOST", "localhost", "PORT", "8080")
	groups := grouper.GroupBy(in, grouper.Options{IncludeUngrouped: false})

	for _, g := range groups {
		if g.Name == "(ungrouped)" {
			t.Error("(ungrouped) group should be excluded")
		}
	}
}

func TestGroupBy_CustomSeparator(t *testing.T) {
	in := entries("DB.HOST", "localhost", "DB.PORT", "5432")
	groups := grouper.GroupBy(in, grouper.Options{Separator: "."})

	if len(groups) != 1 || groups[0].Name != "DB" {
		t.Errorf("expected single group DB, got %+v", groups)
	}
}

func TestGroupBy_EmptyInput(t *testing.T) {
	groups := grouper.GroupBy(nil, grouper.Options{})
	if len(groups) != 0 {
		t.Errorf("expected empty result, got %d groups", len(groups))
	}
}

func TestKeys_ReturnsSortedPrefixes(t *testing.T) {
	in := entries("S3_BUCKET", "x", "APP_NAME", "y", "DB_HOST", "z", "PLAIN", "v")
	keys := grouper.Keys(in, "")

	expected := []string{"APP", "DB", "S3"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, keys)
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], k)
		}
	}
}
