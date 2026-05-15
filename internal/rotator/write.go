package rotator

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// writeEntries serialises entries back to a .env file, preserving KEY=VALUE
// format and quoting values that contain spaces.
func writeEntries(path string, entries []parser.Entry) error {
	var sb strings.Builder
	for _, e := range entries {
		val := e.Value
		if strings.ContainsAny(val, " \t") {
			val = fmt.Sprintf("%q", val)
		}
		fmt.Fprintf(&sb, "%s=%s\n", e.Key, val)
	}
	return os.WriteFile(path, []byte(sb.String()), 0o600)
}
