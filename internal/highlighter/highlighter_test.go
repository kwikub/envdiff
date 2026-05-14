package highlighter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/highlighter"
)

func results(rs ...differ.Result) []differ.Result { return rs }

func noColor() highlighter.Options { return highlighter.Options{NoColor: true} }

func TestHighlight_Added(t *testing.T) {
	out := highlighter.Highlight(results(
		differ.Result{Key: "FOO", Status: differ.Added, NewValue: "bar"},
	), noColor())
	if !strings.Contains(out, "+ FOO=bar") {
		t.Errorf("expected added line, got: %q", out)
	}
}

func TestHighlight_Removed(t *testing.T) {
	out := highlighter.Highlight(results(
		differ.Result{Key: "FOO", Status: differ.Removed, OldValue: "bar"},
	), noColor())
	if !strings.Contains(out, "- FOO=bar") {
		t.Errorf("expected removed line, got: %q", out)
	}
}

func TestHighlight_Modified(t *testing.T) {
	out := highlighter.Highlight(results(
		differ.Result{Key: "FOO", Status: differ.Modified, OldValue: "old", NewValue: "new"},
	), noColor())
	if !strings.Contains(out, "- FOO=old") {
		t.Errorf("expected old value line, got: %q", out)
	}
	if !strings.Contains(out, "+ FOO=new") {
		t.Errorf("expected new value line, got: %q", out)
	}
}

func TestHighlight_Unchanged_HiddenByDefault(t *testing.T) {
	out := highlighter.Highlight(results(
		differ.Result{Key: "FOO", Status: differ.Unchanged, OldValue: "bar"},
	), noColor())
	if strings.Contains(out, "FOO") {
		t.Errorf("expected unchanged key to be hidden, got: %q", out)
	}
}

func TestHighlight_Unchanged_ShownWhenEnabled(t *testing.T) {
	opts := highlighter.Options{NoColor: true, ShowUnchanged: true}
	out := highlighter.Highlight(results(
		differ.Result{Key: "FOO", Status: differ.Unchanged, OldValue: "bar"},
	), opts)
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected unchanged key to appear, got: %q", out)
	}
}

func TestHighlight_Color_ContainsEscapeCodes(t *testing.T) {
	opts := highlighter.Options{NoColor: false}
	out := highlighter.Highlight(results(
		differ.Result{Key: "X", Status: differ.Added, NewValue: "1"},
	), opts)
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI escape codes in output, got: %q", out)
	}
}

func TestSummary_Counts(t *testing.T) {
	rs := results(
		differ.Result{Status: differ.Added},
		differ.Result{Status: differ.Added},
		differ.Result{Status: differ.Removed},
		differ.Result{Status: differ.Modified},
		differ.Result{Status: differ.Unchanged},
	)
	summary := highlighter.Summary(rs, noColor())
	if !strings.Contains(summary, "+2 added") {
		t.Errorf("expected '+2 added' in summary, got: %q", summary)
	}
	if !strings.Contains(summary, "-1 removed") {
		t.Errorf("expected '-1 removed' in summary, got: %q", summary)
	}
	if !strings.Contains(summary, "~1 modified") {
		t.Errorf("expected '~1 modified' in summary, got: %q", summary)
	}
}

func TestHighlight_EmptyResults(t *testing.T) {
	out := highlighter.Highlight(nil, noColor())
	if out != "" {
		t.Errorf("expected empty output for nil results, got: %q", out)
	}
}
