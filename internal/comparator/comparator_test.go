package comparator_test

import (
	"os"
	"testing"

	"github.com/user/envdiff/internal/comparator"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestCompare_NoFiles_ReturnsError(t *testing.T) {
	_, err := comparator.Compare(map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty file map, got nil")
	}
}

func TestCompare_UniformKey(t *testing.T) {
	dev := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	prod := writeTempEnv(t, "APP_ENV=production\nDB_HOST=db.prod.example.com\n")

	report, err := comparator.Compare(map[string]string{"dev": dev, "prod": prod})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, kr := range report.Keys {
		if kr.Key == "APP_ENV" && !kr.Uniform {
			t.Errorf("expected APP_ENV to be uniform across envs")
		}
		if kr.Key == "DB_HOST" && kr.Uniform {
			t.Errorf("expected DB_HOST to differ across envs")
		}
	}
}

func TestCompare_MissingKey_TrackedCorrectly(t *testing.T) {
	dev := writeTempEnv(t, "FEATURE_FLAG=true\nSECRET_KEY=abc\n")
	prod := writeTempEnv(t, "SECRET_KEY=abc\n")

	report, err := comparator.Compare(map[string]string{"dev": dev, "prod": prod})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, kr := range report.Keys {
		if kr.Key == "FEATURE_FLAG" {
			if len(kr.Missing) != 1 || kr.Missing[0] != "prod" {
				t.Errorf("expected FEATURE_FLAG missing in prod, got %v", kr.Missing)
			}
			if kr.Uniform {
				t.Errorf("expected FEATURE_FLAG to be non-uniform due to missing env")
			}
		}
	}
}

func TestCompare_EnvironmentsSorted(t *testing.T) {
	a := writeTempEnv(t, "KEY=1\n")
	b := writeTempEnv(t, "KEY=1\n")
	c := writeTempEnv(t, "KEY=1\n")

	report, err := comparator.Compare(map[string]string{"staging": a, "dev": b, "prod": c})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"dev", "prod", "staging"}
	for i, env := range report.Environments {
		if env != expected[i] {
			t.Errorf("environments[%d] = %q, want %q", i, env, expected[i])
		}
	}
}

func TestCompare_InvalidFile_ReturnsError(t *testing.T) {
	_, err := comparator.Compare(map[string]string{"dev": "/nonexistent/path/.env"})
	if err == nil {
		t.Fatal("expected error for invalid file path, got nil")
	}
}
