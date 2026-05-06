package encryptor_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/encryptor"
)

var testKey16 = []byte("0123456789abcdef")          // 16 bytes
var testKey32 = []byte("0123456789abcdef0123456789abcdef") // 32 bytes

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plaintext := "super-secret-value"
	cipher, err := encryptor.Encrypt(plaintext, testKey16)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	got, err := encryptor.Decrypt(cipher, testKey16)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if got != plaintext {
		t.Errorf("expected %q, got %q", plaintext, got)
	}
}

func TestEncrypt_ProducesBase64(t *testing.T) {
	out, err := encryptor.Encrypt("hello", testKey32)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.ContainsAny(out, "+/=") {
		t.Errorf("expected URL-safe base64 without padding, got %q", out)
	}
}

func TestEncrypt_NonDeterministic(t *testing.T) {
	a, _ := encryptor.Encrypt("value", testKey16)
	b, _ := encryptor.Encrypt("value", testKey16)
	if a == b {
		t.Error("expected two encryptions of the same value to differ (random nonce)")
	}
}

func TestEncrypt_InvalidKey(t *testing.T) {
	_, err := encryptor.Encrypt("val", []byte("short"))
	if err != encryptor.ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestDecrypt_InvalidKey(t *testing.T) {
	_, err := encryptor.Decrypt("anything", []byte("bad"))
	if err != encryptor.ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestDecrypt_MalformedCiphertext(t *testing.T) {
	_, err := encryptor.Decrypt("!!!notbase64!!!", testKey16)
	if err != encryptor.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	enc, _ := encryptor.Encrypt("original", testKey16)
	tampered := enc[:len(enc)-4] + "XXXX"
	_, err := encryptor.Decrypt(tampered, testKey16)
	if err != encryptor.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext on tampered input, got %v", err)
	}
}
