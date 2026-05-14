// Package sanitizer transforms .env entries by applying configurable
// sanitization passes: trimming surrounding whitespace from values,
// normalising keys to UPPER_CASE, stripping surrounding quotes, and
// flagging or rejecting keys that violate standard naming conventions.
//
// Usage:
//
//	res, err := sanitizer.Sanitize(entries, sanitizer.Options{
//		TrimValues:    true,
//		UppercaseKeys: true,
//		StripQuotes:   true,
//	})
//
// Each Result carries the transformed Entry and an optional Warning
// string describing any non-fatal issue detected during sanitization.
package sanitizer
