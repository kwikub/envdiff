// Package patcher provides functionality to apply targeted patches to .env files.
//
// A patch is a set of operations (set or delete) that modify key-value entries
// in an existing .env file without altering unrelated keys. After applying
// patches, the caller can serialise the result back to disk using Format.
//
// Example usage:
//
//	res, err := patcher.Apply(".env", []patcher.Patch{
//		{Op: patcher.OpSet,    Key: "DB_HOST",  Value: "localhost"},
//		{Op: patcher.OpDelete, Key: "LEGACY_KEY"},
//	})
//	if err != nil { ... }
//	os.WriteFile(".env", []byte(patcher.Format(res.Entries)), 0o600)
package patcher
