package reconciler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/reconciler"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestReconcile_MissingKey(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	tgt := writeTempEnv(t, "APP_NAME=myapp\n")

	suggestions, err := reconciler.Reconcile(src, tgt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, s := range suggestions {
		if s.Key == "DB_HOST" && s.Action == reconciler.ActionAdd {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ActionAdd suggestion for DB_HOST")
	}
}

func TestReconcile_ExtraKey(t *testing.T) {
	src := writeTempEnv(t, "APP_NAME=myapp\n")
	tgt := writeTempEnv(t, "APP_NAME=myapp\nOLD_KEY=deprecated\n")

	suggestions, err := reconciler.Reconcile(src, tgt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, s := range suggestions {
		if s.Key == "OLD_KEY" && s.Action == reconciler.ActionRemove {
			found = true
		}
	}
	if !found {
		t.Errorf("expected ActionRemove suggestion for OLD_KEY")
	}
}

func TestReconcile_ModifiedKey(t *testing.T) {
	src := writeTempEnv(t, "DB_HOST=prod-db\n")
	tgt := writeTempEnv(t, "DB_HOST=local-db\n")

	suggestions, err := reconciler.Reconcile(src, tgt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, s := range suggestions {
		if s.Key == "DB_HOST" {
			if s.Action != reconciler.ActionUpdate {
				t.Errorf("expected ActionUpdate, got %s", s.Action)
			}
			return
		}
	}
	t.Errorf("DB_HOST suggestion not found")
}

func TestFormat_MasksSecrets(t *testing.T) {
	suggestions := []reconciler.Suggestion{
		{Key: "API_SECRET", Action: reconciler.ActionKeep, Value: "supersecret", Reason: "in sync"},
		{Key: "APP_NAME", Action: reconciler.ActionKeep, Value: "myapp", Reason: "in sync"},
	}

	out := reconciler.Format(suggestions, true)
	if contains(out, "supersecret") {
		t.Errorf("expected secret to be masked, got: %s", out)
	}
	if !contains(out, "myapp") {
		t.Errorf("expected non-secret value to be visible")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
