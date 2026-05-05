package differ

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func TestIsSecret_KnownSecretKeys(t *testing.T) {
	secretKeys := []string{
		"SECRET",
		"API_KEY",
		"DB_PASSWORD",
		"AUTH_TOKEN",
		"PRIVATE_KEY",
		"ACCESS_SECRET",
		"MY_PASSWD",
	}

	for _, key := range secretKeys {
		if !IsSecret(key) {
			t.Errorf("expected IsSecret(%q) = true, got false", key)
		}
	}
}

func TestIsSecret_NonSecretKeys(t *testing.T) {
	nonSecretKeys := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_HOST",
		"FEATURE_FLAG",
	}

	for _, key := range nonSecretKeys {
		if IsSecret(key) {
			t.Errorf("expected IsSecret(%q) = false, got true", key)
		}
	}
}

func TestMaskSecrets_MasksSecretValues(t *testing.T) {
	entries := []parser.Entry{
		{Key: "API_KEY", Value: "super-secret-123"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_PASSWORD", Value: "hunter2"},
		{Key: "PORT", Value: "8080"},
	}

	masked := MaskSecrets(entries)

	expected := map[string]string{
		"API_KEY":     "***",
		"APP_ENV":     "production",
		"DB_PASSWORD": "***",
		"PORT":        "8080",
	}

	for _, e := range masked {
		if want, ok := expected[e.Key]; ok {
			if e.Value != want {
				t.Errorf("key %q: expected value %q, got %q", e.Key, want, e.Value)
			}
		}
	}
}

func TestMaskSecrets_DoesNotMutateOriginal(t *testing.T) {
	original := []parser.Entry{
		{Key: "API_KEY", Value: "my-secret"},
	}

	_ = MaskSecrets(original)

	if original[0].Value != "my-secret" {
		t.Errorf("MaskSecrets mutated original entry: got %q", original[0].Value)
	}
}
