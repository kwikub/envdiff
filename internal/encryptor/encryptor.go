// Package encryptor provides AES-GCM encryption and decryption
// for secret values found in .env files.
package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidKey is returned when the key length is not 16, 24, or 32 bytes.
var ErrInvalidKey = errors.New("encryptor: key must be 16, 24, or 32 bytes")

// ErrInvalidCiphertext is returned when decryption fails due to malformed input.
var ErrInvalidCiphertext = errors.New("encryptor: invalid ciphertext")

// Encrypt encrypts plaintext using AES-GCM with the provided key.
// The result is base64-encoded (URL-safe, no padding).
func Encrypt(plaintext string, key []byte) (string, error) {
	if err := validateKey(key); err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded AES-GCM ciphertext using the provided key.
func Decrypt(encoded string, key []byte) (string, error) {
	if err := validateKey(key); err != nil {
		return "", err
	}
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", ErrInvalidCiphertext
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	return string(plaintext), nil
}

func validateKey(key []byte) error {
	switch len(key) {
	case 16, 24, 32:
		return nil
	}
	return ErrInvalidKey
}
