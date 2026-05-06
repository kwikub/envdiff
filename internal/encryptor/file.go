package encryptor

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// EncryptedEntry holds a key with its encrypted value.
type EncryptedEntry struct {
	Key   string
	Value string // base64-encoded ciphertext
}

// EncryptFile reads a parsed .env file and encrypts values for keys that
// match the provided key set (nil means encrypt all values).
func EncryptFile(entries []parser.Entry, key []byte, targets map[string]bool) ([]EncryptedEntry, error) {
	out := make([]EncryptedEntry, 0, len(entries))
	for _, e := range entries {
		val := e.Value
		if targets == nil || targets[e.Key] {
			encVal, err := Encrypt(e.Value, key)
			if err != nil {
				return nil, fmt.Errorf("encrypt key %q: %w", e.Key, err)
			}
			val = encVal
		}
		out = append(out, EncryptedEntry{Key: e.Key, Value: val})
	}
	return out, nil
}

// DecryptFile decrypts previously encrypted entries back to plaintext.
// Keys not present in targets (when non-nil) are passed through unchanged.
func DecryptFile(entries []EncryptedEntry, key []byte, targets map[string]bool) ([]parser.Entry, error) {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		val := e.Value
		if targets == nil || targets[e.Key] {
			decVal, err := Decrypt(e.Value, key)
			if err != nil {
				return nil, fmt.Errorf("decrypt key %q: %w", e.Key, err)
			}
			val = decVal
		}
		out = append(out, parser.Entry{Key: e.Key, Value: val})
	}
	return out, nil
}
