package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/reporter"
)

func TestReportText_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	if err := reporter.Report(&buf, nil, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %q", buf.String())
	}
}

func TestReportText_Added(t *testing.T) {
	diffs := []differ.DiffEntry{
		{Key: "NEW_KEY", Status: differ.StatusAdded, NewValue: "hello"},
	}
	var buf bytes.Buffer
	if err := reporter.Report(&buf, diffs, reporter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "+ NEW_KEY=hello") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestReportText_Removed(t *testing.T) {
	diffs := []differ.DiffEntry{
		{Key: "OLD_KEY", Status: differ.StatusRemoved, OldValue: "bye"},
	}
	var buf bytes.Buffer
	reporter.Report(&buf, diffs, reporter.FormatText)
	if !strings.Contains(buf.String(), "- OLD_KEY=bye") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestReportText_Modified(t *testing.T) {
	diffs := []differ.DiffEntry{
		{Key: "PORT", Status: differ.StatusModified, OldValue: "3000", NewValue: "4000"},
	}
	var buf bytes.Buffer
	reporter.Report(&buf, diffs, reporter.FormatText)
	if !strings.Contains(buf.String(), "~ PORT: 3000 -> 4000") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestReportJSON_ContainsKey(t *testing.T) {
	diffs := []differ.DiffEntry{
		{Key: "API_KEY", Status: differ.StatusAdded, NewValue: "secret"},
	}
	var buf bytes.Buffer
	if err := reporter.Report(&buf, diffs, reporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"API_KEY"`) {
		t.Errorf("expected key in JSON output, got: %q", out)
	}
	if !strings.Contains(out, `"added"`) {
		t.Errorf("expected status in JSON output, got: %q", out)
	}
}
