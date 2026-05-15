package cloner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/cloner"
	"github.com/user/envdiff/internal/parser"
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

func readKeys(t *testing.T, path string) map[string]string {
	t.Helper()
	entries, err := parser.Parse(path)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	m := map[string]string{}
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

func TestClone_BasicCopy(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	dst := writeTempEnv(t, "EXISTING=yes\n")

	res, err := cloner.Clone(src, dst, cloner.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	m := readKeys(t, dst)
	if m["FOO"] != "bar" || m["BAZ"] != "qux" || m["EXISTING"] != "yes" {
		t.Errorf("unexpected destination state: %v", m)
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := writeTempEnv(t, "FOO=new\n")
	dst := writeTempEnv(t, "FOO=old\n")

	res, err := cloner.Clone(src, dst, cloner.Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected FOO to be skipped, got %v", res.Skipped)
	}
	m := readKeys(t, dst)
	if m["FOO"] != "old" {
		t.Errorf("expected old value preserved, got %q", m["FOO"])
	}
}

func TestClone_OverwriteExisting(t *testing.T) {
	src := writeTempEnv(t, "FOO=new\n")
	dst := writeTempEnv(t, "FOO=old\n")

	_, err := cloner.Clone(src, dst, cloner.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := readKeys(t, dst)
	if m["FOO"] != "new" {
		t.Errorf("expected new value, got %q", m["FOO"])
	}
}

func TestClone_PrefixFilter(t *testing.T) {
	src := writeTempEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_URL=postgres\n")
	dst := writeTempEnv(t, "")

	res, err := cloner.Clone(src, dst, cloner.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	m := readKeys(t, dst)
	if _, ok := m["DB_URL"]; ok {
		t.Error("DB_URL should not have been copied")
	}
}

func TestClone_AllowList(t *testing.T) {
	src := writeTempEnv(t, "FOO=1\nBAR=2\nBAZ=3\n")
	dst := writeTempEnv(t, "")

	res, err := cloner.Clone(src, dst, cloner.Options{AllowList: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 2 {
		t.Errorf("expected 2 copied, got %d", len(res.Copied))
	}
	m := readKeys(t, dst)
	if _, ok := m["BAR"]; ok {
		t.Error("BAR should not have been copied")
	}
}

func TestClone_DryRun_DoesNotModifyFile(t *testing.T) {
	src := writeTempEnv(t, "FOO=bar\n")
	dst := writeTempEnv(t, "EXISTING=yes\n")

	before, _ := os.ReadFile(dst)
	res, err := cloner.Clone(src, dst, cloner.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Copied) != 1 {
		t.Errorf("expected 1 in dry-run result, got %d", len(res.Copied))
	}
	after, _ := os.ReadFile(dst)
	if string(before) != string(after) {
		t.Error("dry-run should not modify the destination file")
	}
}

func TestClone_InvalidSource_ReturnsError(t *testing.T) {
	dst := writeTempEnv(t, "")
	_, err := cloner.Clone(filepath.Join(t.TempDir(), "missing.env"), dst, cloner.Options{})
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
