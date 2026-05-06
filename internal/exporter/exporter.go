// Package exporter provides functionality to export reconciled .env files
// in various formats (shell export, Docker env-file, JSON).
package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Format represents the output format for export.
type Format string

const (
	FormatEnv    Format = "env"
	FormatShell  Format = "shell"
	FormatDocker Format = "docker"
	FormatJSON   Format = "json"
)

// Export writes entries to w in the specified format.
func Export(entries []parser.Entry, format Format, w io.Writer) error {
	switch format {
	case FormatEnv:
		return exportEnv(entries, w)
	case FormatShell:
		return exportShell(entries, w)
	case FormatDocker:
		return exportDocker(entries, w)
	case FormatJSON:
		return exportJSON(entries, w)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func exportEnv(entries []parser.Entry, w io.Writer) error {
	for _, e := range entries {
		val := quoteIfNeeded(e.Value)
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func exportShell(entries []parser.Entry, w io.Writer) error {
	for _, e := range entries {
		val := quoteIfNeeded(e.Value)
		if _, err := fmt.Fprintf(w, "export %s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func exportDocker(entries []parser.Entry, w io.Writer) error {
	for _, e := range entries {
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value); err != nil {
			return err
		}
	}
	return nil
}

func exportJSON(entries []parser.Entry, w io.Writer) error {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}

func quoteIfNeeded(val string) string {
	if strings.ContainsAny(val, " \t#") {
		return fmt.Sprintf("%q", val)
	}
	return val
}
