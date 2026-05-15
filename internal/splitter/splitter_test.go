package splitter_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/splitter"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestSplit_ByPrefix_CreatesCorrectFiles(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\nAPP_ENV=prod\n")
	dest := t.TempDir()

	counts, err := splitter.Split(src, dest, splitter.Options{Strategy: splitter.StrategyPrefix})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if counts["db.env"] != 2 {
		t.Errorf("expected 2 DB entries, got %d", counts["db.env"])
	}
	if counts["app.env"] != 2 {
		t.Errorf("expected 2 APP entries, got %d", counts["app.env"])
	}
}

func TestSplit_ByPrefix_FileContainsCorrectKeys(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	dest := t.TempDir()

	_, err := splitter.Split(src, dest, splitter.Options{Strategy: splitter.StrategyPrefix})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dest, "db.env"))
	if err != nil {
		t.Fatalf("could not read db.env: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in db.env, got:\n%s", content)
	}
}

func TestSplit_AlphaStrategy_ProducesBuckets(t *testing.T) {
	src := writeTempEnv(t, "ALPHA=1\nBETA=2\nGAMMA=3\nDELTA=4\n")
	dest := t.TempDir()

	counts, err := splitter.Split(src, dest, splitter.Options{
		Strategy: splitter.StrategyAlpha,
		Buckets:  4,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	total := 0
	for _, c := range counts {
		total += c
	}
	if total != 4 {
		t.Errorf("expected 4 total entries across buckets, got %d", total)
	}
}

func TestSplit_CustomMapping_AssignsCorrectly(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=localhost\nREDIS_URL=redis://x\nAPP_NAME=myapp\n")
	dest := t.TempDir()

	counts, err := splitter.Split(src, dest, splitter.Options{
		Mapping: map[string][]string{
			"databases.env": {"DB_", "REDIS_"},
			"app.env":       {"APP_"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if counts["databases.env"] != 2 {
		t.Errorf("expected 2 database entries, got %d", counts["databases.env"])
	}
	if counts["app.env"] != 1 {
		t.Errorf("expected 1 app entry, got %d", counts["app.env"])
	}
}

func TestSplit_CustomMapping_UnmatchedGoesToOther(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=localhost\nUNKNOWN_KEY=value\n")
	dest := t.TempDir()

	counts, err := splitter.Split(src, dest, splitter.Options{
		Mapping: map[string][]string{
			"db.env": {"DB_"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if counts["other.env"] != 1 {
		t.Errorf("expected 1 entry in other.env, got %d", counts["other.env"])
	}
}
