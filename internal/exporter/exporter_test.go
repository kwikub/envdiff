package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
)

var testEntries = []parser.Entry{
	{Key: "APP_NAME", Value: "myapp"},
	{Key: "PORT", Value: "8080"},
	{Key: "DB_URL", Value: "postgres://localhost/db"},
	{Key: "GREETING", Value: "hello world"},
}

func TestExport_EnvFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(testEntries, exporter.FormatEnv, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME=myapp in output, got:\n%s", out)
	}
	if !strings.Contains(out, `GREETING="hello world"`) {
		t.Errorf("expected quoted GREETING in output, got:\n%s", out)
	}
}

func TestExport_ShellFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(testEntries, exporter.FormatShell, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_NAME=myapp") {
		t.Errorf("expected 'export APP_NAME=myapp' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "export PORT=8080") {
		t.Errorf("expected 'export PORT=8080' in output, got:\n%s", out)
	}
}

func TestExport_DockerFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(testEntries, exporter.FormatDocker, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	// Docker format should not quote values
	if !strings.Contains(out, "GREETING=hello world") {
		t.Errorf("expected unquoted GREETING in docker format, got:\n%s", out)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.Export(testEntries, exporter.FormatJSON, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME=myapp, got %s", m["APP_NAME"])
	}
	if m["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", m["PORT"])
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(testEntries, exporter.Format("xml"), &buf)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
