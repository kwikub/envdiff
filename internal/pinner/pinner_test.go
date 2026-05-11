package pinner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/pinner"
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

func TestPin_CreatesLockfile(t *testing.T) {
	src := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	dest := filepath.Join(t.TempDir(), "env.lock")

	pf, err := pinner.Pin(src, dest)
	if err != nil {
		t.Fatalf("Pin() error: %v", err)
	}
	if pf.Entries["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", pf.Entries["APP_ENV"])
	}
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("lockfile not created: %v", err)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	src := writeTempEnv(t, "KEY=value\n")
	dest := filepath.Join(t.TempDir(), "env.lock")

	if _, err := pinner.Pin(src, dest); err != nil {
		t.Fatal(err)
	}
	pf, err := pinner.Load(dest)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if pf.Entries["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", pf.Entries["KEY"])
	}
}

func TestVerify_NoDrift(t *testing.T) {
	src := writeTempEnv(t, "APP_ENV=staging\nPORT=8080\n")
	dest := filepath.Join(t.TempDir(), "env.lock")

	if _, err := pinner.Pin(src, dest); err != nil {
		t.Fatal(err)
	}
	drifted, err := pinner.Verify(src, dest)
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}
	if len(drifted) != 0 {
		t.Errorf("expected no drift, got %v", drifted)
	}
}

func TestVerify_DetectsDrift(t *testing.T) {
	original := writeTempEnv(t, "APP_ENV=staging\nPORT=8080\n")
	dest := filepath.Join(t.TempDir(), "env.lock")

	if _, err := pinner.Pin(original, dest); err != nil {
		t.Fatal(err)
	}

	modified := writeTempEnv(t, "APP_ENV=production\nPORT=8080\n")
	drifted, err := pinner.Verify(modified, dest)
	if err != nil {
		t.Fatalf("Verify() error: %v", err)
	}
	if len(drifted) != 1 || drifted[0] != "APP_ENV" {
		t.Errorf("expected [APP_ENV] drift, got %v", drifted)
	}
}

func TestLoad_InvalidFile_ReturnsError(t *testing.T) {
	_, err := pinner.Load("/nonexistent/env.lock")
	if err == nil {
		t.Error("expected error for missing lockfile")
	}
}
