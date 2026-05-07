// Package inspector analyses parsed .env files and produces a Summary
// containing aggregate metadata such as total key count, number of secret
// keys, empty values, duplicate keys, and key groupings derived from
// underscore-separated prefixes.
//
// Usage:
//
//	s, err := inspector.Inspect(".env")
//	if err != nil { ... }
//	fmt.Println(s.TotalKeys, s.SecretKeys)
package inspector
