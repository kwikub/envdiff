package exporter

import "fmt"

// ParseFormat converts a string to a Format, returning an error if unrecognised.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatEnv, FormatShell, FormatDocker, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf(
			"unknown format %q: must be one of env, shell, docker, json", s,
		)
	}
}

// String implements the Stringer interface for Format.
func (f Format) String() string {
	return string(f)
}

// ValidFormats returns all supported export format strings.
func ValidFormats() []string {
	return []string{
		string(FormatEnv),
		string(FormatShell),
		string(FormatDocker),
		string(FormatJSON),
	}
}
