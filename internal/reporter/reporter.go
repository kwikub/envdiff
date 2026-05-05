// Package reporter formats and outputs diff results for human or machine consumption.
package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Format controls the output format of the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a formatted diff report to w.
func Report(w io.Writer, diffs []differ.DiffEntry, format Format) error {
	switch format {
	case FormatJSON:
		return reportJSON(w, diffs)
	default:
		return reportText(w, diffs)
	}
}

func reportText(w io.Writer, diffs []differ.DiffEntry) error {
	if len(diffs) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}
	for _, d := range diffs {
		var line string
		switch d.Status {
		case differ.StatusAdded:
			line = fmt.Sprintf("+ %s=%s", d.Key, d.NewValue)
		case differ.StatusRemoved:
			line = fmt.Sprintf("- %s=%s", d.Key, d.OldValue)
		case differ.StatusModified:
			line = fmt.Sprintf("~ %s: %s -> %s", d.Key, d.OldValue, d.NewValue)
		case differ.StatusUnchanged:
			line = fmt.Sprintf("  %s=%s", d.Key, d.NewValue)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func reportJSON(w io.Writer, diffs []differ.DiffEntry) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, d := range diffs {
		sb.WriteString(fmt.Sprintf(
			"  {\"key\": %q, \"status\": %q, \"old_value\": %q, \"new_value\": %q}",
			d.Key, d.Status, d.OldValue, d.NewValue,
		))
		if i < len(diffs)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}
