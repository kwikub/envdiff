package encryptor_test

import (
	"testing"

	"github.com/user/envdiff/internal/encryptor"
	"github.com/user/envdiff/internal/parser"
)

func TestEncryptFile_AllKeys(t *testing.T) {
	entries := []parser.Entry{
		{Key: "DB_PASS", Value: "secret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	enc, err := encryptor.EncryptFile(entries, testKey16, nil)
	if err != nil {
		t.Fatalf("EncryptFile error: %v", err)
	}
	if len(enc) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(enc))
	}
	for i, e := range enc {
		if e.Value == entries[i].Value {
			t.Errorf("key %q: value was not encrypted", e.Key)
		}
	}
}

func TestEncryptFile_TargetedKeys(t *testing.T) {
	entries := []parser.Entry{
		{Key: "DB_PASS", Value: "secret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	targets := map[string]bool{"DB_PASS": true}
	enc, err := encryptor.EncryptFile(entries, testKey16, targets)
	if err != nil {
		t.Fatalf("EncryptFile error: %v", err)
	}
	if enc[0].Value == "secret" {
		t.Error("DB_PASS should have been encrypted")
	}
	if enc[1].Value != "myapp" {
		t.Errorf("APP_NAME should be unchanged, got %q", enc[1].Value)
	}
}

func TestDecryptFile_RoundTrip(t *testing.T) {
	original := []parser.Entry{
		{Key: "SECRET_KEY", Value: "abc123"},
		{Key: "HOST", Value: "localhost"},
	}
	enc, err := encryptor.EncryptFile(original, testKey32, nil)
	if err != nil {
		t.Fatalf("EncryptFile: %v", err)
	}
	dec, err := encryptor.DecryptFile(enc, testKey32, nil)
	if err != nil {
		t.Fatalf("DecryptFile: %v", err)
	}
	for i, e := range dec {
		if e.Key != original[i].Key || e.Value != original[i].Value {
			t.Errorf("entry %d: expected %+v, got %+v", i, original[i], e)
		}
	}
}

func TestDecryptFile_WrongKey(t *testing.T) {
	original := []parser.Entry{{Key: "TOKEN", Value: "tok"}}
	enc, _ := encryptor.EncryptFile(original, testKey16, nil)
	wrongKey := []byte("fedcba9876543210")
	_, err := encryptor.DecryptFile(enc, wrongKey, nil)
	if err == nil {
		t.Error("expected error when decrypting with wrong key")
	}
}
