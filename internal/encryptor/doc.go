// Package encryptor provides AES-GCM encryption and decryption utilities
// for protecting secret values stored in .env files.
//
// # Overview
//
// Use [Encrypt] and [Decrypt] for individual string values, and
// [EncryptFile] / [DecryptFile] for operating on a full slice of
// parsed entries (as returned by the parser package).
//
// Keys must be exactly 16, 24, or 32 bytes to select AES-128, AES-192,
// or AES-256 respectively. All ciphertexts are encoded as URL-safe
// base64 without padding so they are safe to embed directly in .env files.
package encryptor
