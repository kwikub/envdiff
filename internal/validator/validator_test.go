package validator_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestValidate_AllKeysPresent(t *testing.T) {
	tmpl := writeTempEnv(t, "APP_ENV=development\nDB_HOST=localhost\n")
	env := writeTempEnv(t, "APP_ENV=production\nDB_HOST=db.example.com\n")

	res, err := validator.Validate(env, tmpl, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Valid {
		t.Errorf("expected valid result, got errors: %v", res.Errors)
	}
}

func TestValidate_MissingKey(t *testing.T) {
	tmpl := writeTempEnv(t, "APP_ENV=development\nDB_HOST=localhost\n")
	env := writeTempEnv(t, "APP_ENV=production\n")

	res, err := validator.Validate(env, tmpl, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Valid {
		t.Error("expected invalid result due to missing key")
	}
	if len(res.Errors) != 1 || res.Errors[0].Key != "DB_HOST" {
		t.Errorf("expected error for DB_HOST, got: %v", res.Errors)
	}
}

func TestValidate_PatternRule_Pass(t *testing.T) {
	tmpl := writeTempEnv(t, "PORT=8080\n")
	env := writeTempEnv(t, "PORT=3000\n")

	rules := []validator.Rule{
		{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
	}
	res, err := validator.Validate(env, tmpl, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Valid {
		t.Errorf("expected valid result, got: %v", res.Errors)
	}
}

func TestValidate_PatternRule_Fail(t *testing.T) {
	tmpl := writeTempEnv(t, "PORT=8080\n")
	env := writeTempEnv(t, "PORT=not-a-number\n")

	rules := []validator.Rule{
		{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
	}
	res, err := validator.Validate(env, tmpl, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Valid {
		t.Error("expected invalid result due to pattern mismatch")
	}
	if len(res.Errors) == 0 || res.Errors[0].Key != "PORT" {
		t.Errorf("expected PORT pattern error, got: %v", res.Errors)
	}
}

func TestValidate_RequiredRuleNotInTemplate(t *testing.T) {
	tmpl := writeTempEnv(t, "APP_ENV=dev\n")
	env := writeTempEnv(t, "APP_ENV=prod\n")

	rules := []validator.Rule{
		{Key: "EXTRA_REQUIRED", Required: true},
	}
	res, err := validator.Validate(env, tmpl, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Valid {
		t.Error("expected invalid result for missing required key not in template")
	}
}
